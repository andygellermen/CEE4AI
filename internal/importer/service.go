package importer

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	pool *pgxpool.Pool
}

type columnType struct {
	DataType string
	UDTName  string
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

func (s *Service) SeedDir(ctx context.Context, root string) error {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin seed transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	for _, spec := range DefaultTableSpecs() {
		path := filepath.Join(root, spec.RelativePath)
		csvData, err := ReadCSV(path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return fmt.Errorf("read seed file %s: %w", path, err)
		}

		if len(csvData.Rows) == 0 {
			continue
		}

		if err := s.importRows(ctx, tx, spec, csvData); err != nil {
			return fmt.Errorf("import %s: %w", spec.TableName, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit seed transaction: %w", err)
	}

	return nil
}

func (s *Service) importRows(ctx context.Context, tx pgx.Tx, spec TableSpec, csvData *CSVData) error {
	schemaName, tableName, err := splitTableName(spec.TableName)
	if err != nil {
		return err
	}

	columnTypes, err := loadColumnTypes(ctx, tx, schemaName, tableName)
	if err != nil {
		return err
	}

	query := buildUpsertQuery(spec.TableName, csvData.Headers, spec.ConflictColumns)
	for _, row := range csvData.Rows {
		values := make([]any, 0, len(csvData.Headers))
		for _, header := range csvData.Headers {
			column, ok := columnTypes[header]
			if !ok {
				return fmt.Errorf("column %s does not exist on %s", header, spec.TableName)
			}

			parsed, err := parseValue(column, row[header])
			if err != nil {
				return fmt.Errorf("parse %s.%s value %q: %w", spec.TableName, header, row[header], err)
			}
			values = append(values, parsed)
		}

		if _, err := tx.Exec(ctx, query, values...); err != nil {
			return fmt.Errorf("exec upsert on %s: %w", spec.TableName, err)
		}
	}

	if spec.ResetSequence && contains(csvData.Headers, "id") {
		if err := syncSequence(ctx, tx, spec.TableName); err != nil {
			return err
		}
	}

	return nil
}

func loadColumnTypes(ctx context.Context, tx pgx.Tx, schemaName, tableName string) (map[string]columnType, error) {
	rows, err := tx.Query(ctx, `
SELECT column_name, data_type, udt_name
FROM information_schema.columns
WHERE table_schema = $1 AND table_name = $2
ORDER BY ordinal_position
`, schemaName, tableName)
	if err != nil {
		return nil, fmt.Errorf("query column types for %s.%s: %w", schemaName, tableName, err)
	}
	defer rows.Close()

	columns := make(map[string]columnType)
	for rows.Next() {
		var name string
		var dataType string
		var udtName string
		if err := rows.Scan(&name, &dataType, &udtName); err != nil {
			return nil, fmt.Errorf("scan column types: %w", err)
		}
		columns[name] = columnType{DataType: dataType, UDTName: udtName}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate column types: %w", err)
	}

	return columns, nil
}

func parseValue(column columnType, raw string) (any, error) {
	if raw == "" {
		return nil, nil
	}

	switch column.DataType {
	case "boolean":
		return strconv.ParseBool(raw)
	case "smallint", "integer", "bigint":
		return strconv.ParseInt(raw, 10, 64)
	case "numeric", "real", "double precision":
		return strconv.ParseFloat(raw, 64)
	default:
		switch column.UDTName {
		case "bool":
			return strconv.ParseBool(raw)
		case "int2", "int4", "int8":
			return strconv.ParseInt(raw, 10, 64)
		case "float4", "float8", "numeric":
			return strconv.ParseFloat(raw, 64)
		default:
			return raw, nil
		}
	}
}

func buildUpsertQuery(tableName string, headers, conflictColumns []string) string {
	quotedHeaders := make([]string, 0, len(headers))
	placeholders := make([]string, 0, len(headers))
	updates := make([]string, 0, len(headers))

	for idx, header := range headers {
		quotedHeaders = append(quotedHeaders, quoteIdentifier(header))
		placeholders = append(placeholders, fmt.Sprintf("$%d", idx+1))
		if contains(conflictColumns, header) {
			continue
		}
		updates = append(updates, fmt.Sprintf("%s = EXCLUDED.%s", quoteIdentifier(header), quoteIdentifier(header)))
	}

	conflicts := make([]string, 0, len(conflictColumns))
	for _, column := range conflictColumns {
		conflicts = append(conflicts, quoteIdentifier(column))
	}

	base := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) ON CONFLICT (%s)",
		tableName,
		strings.Join(quotedHeaders, ", "),
		strings.Join(placeholders, ", "),
		strings.Join(conflicts, ", "),
	)

	if len(updates) == 0 {
		return base + " DO NOTHING"
	}

	return base + " DO UPDATE SET " + strings.Join(updates, ", ")
}

func syncSequence(ctx context.Context, tx pgx.Tx, tableName string) error {
	var sequenceName *string
	if err := tx.QueryRow(ctx, `SELECT pg_get_serial_sequence($1, 'id')`, tableName).Scan(&sequenceName); err != nil {
		return fmt.Errorf("fetch sequence for %s: %w", tableName, err)
	}

	if sequenceName == nil || *sequenceName == "" {
		return nil
	}

	query := fmt.Sprintf(
		"SELECT setval(%s, COALESCE((SELECT MAX(id) FROM %s), 1), true)",
		quoteLiteral(*sequenceName),
		tableName,
	)

	if _, err := tx.Exec(ctx, query); err != nil {
		return fmt.Errorf("sync sequence for %s: %w", tableName, err)
	}

	return nil
}

func splitTableName(tableName string) (string, string, error) {
	parts := strings.Split(tableName, ".")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("table name %q must be schema-qualified", tableName)
	}
	return parts[0], parts[1], nil
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func quoteIdentifier(value string) string {
	return `"` + strings.ReplaceAll(value, `"`, `""`) + `"`
}

func quoteLiteral(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "''") + "'"
}
