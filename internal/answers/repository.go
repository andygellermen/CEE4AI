package answers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrAnswerNotFound = errors.New("answer not found")

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, params CreateAnswerParams) (*Answer, error) {
	row := r.pool.QueryRow(ctx, `
INSERT INTO runtime.answers (
    session_id,
    package_id,
    question_id,
    answer_kind,
    selected_option_ids,
    scale_value,
    free_text_answer,
    raw_score,
    evaluated_score,
    certainty_level
)
VALUES ($1::uuid, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING
    id,
    session_id::text,
    package_id,
    question_id,
    answer_kind,
    selected_option_ids,
    scale_value,
    free_text_answer,
    raw_score,
    evaluated_score,
    certainty_level,
    answered_at,
    created_at
`, params.SessionID, params.PackageID, params.QuestionID, params.AnswerKind, params.SelectedOptionIDs, params.ScaleValue, nullIfEmpty(params.FreeTextAnswer), params.RawScore, params.EvaluatedScore, nullIfEmpty(params.CertaintyLevel))

	return scanAnswer(row)
}

func (r *Repository) ExistsForSessionQuestion(ctx context.Context, sessionID string, questionID int64) (bool, error) {
	var exists bool
	if err := r.pool.QueryRow(ctx, `
SELECT EXISTS(
    SELECT 1
    FROM runtime.answers
    WHERE session_id = $1::uuid
      AND question_id = $2
)
`, sessionID, questionID).Scan(&exists); err != nil {
		return false, fmt.Errorf("query answer existence: %w", err)
	}

	return exists, nil
}

func (r *Repository) CountForSession(ctx context.Context, sessionID string) (int, error) {
	var count int
	if err := r.pool.QueryRow(ctx, `
SELECT COUNT(*)
FROM runtime.answers
WHERE session_id = $1::uuid
`, sessionID).Scan(&count); err != nil {
		return 0, fmt.Errorf("count answers: %w", err)
	}

	return count, nil
}

func scanAnswer(row pgx.Row) (*Answer, error) {
	var answer Answer
	var packageID sql.NullInt64
	var selectedOptionIDs []byte
	var scaleValue sql.NullInt64
	var freeText sql.NullString
	var rawScore sql.NullFloat64
	var evaluatedScore sql.NullFloat64
	var certainty sql.NullString

	if err := row.Scan(
		&answer.ID,
		&answer.SessionID,
		&packageID,
		&answer.QuestionID,
		&answer.AnswerKind,
		&selectedOptionIDs,
		&scaleValue,
		&freeText,
		&rawScore,
		&evaluatedScore,
		&certainty,
		&answer.AnsweredAt,
		&answer.CreatedAt,
	); err != nil {
		return nil, err
	}

	if packageID.Valid {
		value := packageID.Int64
		answer.PackageID = &value
	}
	if len(selectedOptionIDs) > 0 {
		if err := json.Unmarshal(selectedOptionIDs, &answer.SelectedOptionIDs); err != nil {
			return nil, fmt.Errorf("unmarshal selected option ids: %w", err)
		}
	}
	if scaleValue.Valid {
		value := int(scaleValue.Int64)
		answer.ScaleValue = &value
	}
	if freeText.Valid {
		answer.FreeTextAnswer = freeText.String
	}
	if rawScore.Valid {
		value := rawScore.Float64
		answer.RawScore = &value
	}
	if evaluatedScore.Valid {
		value := evaluatedScore.Float64
		answer.EvaluatedScore = &value
	}
	if certainty.Valid {
		answer.CertaintyLevel = certainty.String
	}

	return &answer, nil
}

func nullIfEmpty(value string) any {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return strings.TrimSpace(value)
}
