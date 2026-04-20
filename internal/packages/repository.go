package packages

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrPackageNotFound = errors.New("package not found")

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) GetBySessionAndIndex(ctx context.Context, sessionID string, packageIndex int) (*SessionPackage, error) {
	row := r.pool.QueryRow(ctx, `
SELECT
    id,
    session_id::text,
    package_index,
    package_size,
    estimated_time_seconds,
    actual_time_seconds,
    completion_quality,
    continuation_window_until,
    recommended_next_mode,
    created_at
FROM runtime.session_packages
WHERE session_id = $1::uuid AND package_index = $2
`, sessionID, packageIndex)

	pkg, err := scanPackage(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrPackageNotFound
	}

	return pkg, err
}

func (r *Repository) Create(ctx context.Context, params CreatePackageParams) (*SessionPackage, error) {
	row := r.pool.QueryRow(ctx, `
INSERT INTO runtime.session_packages (
    session_id,
    package_index,
    package_size,
    estimated_time_seconds,
    recommended_next_mode,
    continuation_window_until
)
VALUES ($1::uuid, $2, $3, $4, $5, $6)
RETURNING
    id,
    session_id::text,
    package_index,
    package_size,
    estimated_time_seconds,
    actual_time_seconds,
    completion_quality,
    continuation_window_until,
    recommended_next_mode,
    created_at
`, params.SessionID, params.PackageIndex, params.PackageSize, params.EstimatedTimeSeconds, params.RecommendedNextMode, params.ContinuationWindowUntil)

	return scanPackage(row)
}

func scanPackage(row pgx.Row) (*SessionPackage, error) {
	var pkg SessionPackage
	var estimated sql.NullInt64
	var actual sql.NullInt64
	var completion sql.NullFloat64
	var continuation sql.NullTime
	var nextMode sql.NullString

	if err := row.Scan(
		&pkg.ID,
		&pkg.SessionID,
		&pkg.PackageIndex,
		&pkg.PackageSize,
		&estimated,
		&actual,
		&completion,
		&continuation,
		&nextMode,
		&pkg.CreatedAt,
	); err != nil {
		return nil, err
	}

	if estimated.Valid {
		value := int(estimated.Int64)
		pkg.EstimatedTimeSeconds = &value
	}
	if actual.Valid {
		value := int(actual.Int64)
		pkg.ActualTimeSeconds = &value
	}
	if completion.Valid {
		value := completion.Float64
		pkg.CompletionQuality = &value
	}
	if continuation.Valid {
		value := continuation.Time
		pkg.ContinuationWindowUntil = &value
	}
	if nextMode.Valid {
		value := nextMode.String
		pkg.RecommendedNextMode = &value
	}

	return &pkg, nil
}
