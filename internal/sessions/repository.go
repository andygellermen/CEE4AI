package sessions

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrSessionNotFound = errors.New("session not found")

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, params CreateSessionParams) (*Session, error) {
	row := r.pool.QueryRow(ctx, `
INSERT INTO runtime.sessions (
    domain_id,
    mode,
    session_goal,
    locale_language_id,
    locale_region_id,
    progress_state
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING
    id::text,
    domain_id,
    mode,
    session_goal,
    locale_language_id,
    locale_region_id,
    result_confidence,
    progress_state,
    started_at,
    finished_at,
    created_at,
    updated_at
`, params.DomainID, params.Mode, nullIfEmpty(params.SessionGoal), params.LocaleLanguageID, params.LocaleRegionID, params.ProgressState)

	return scanSession(row)
}

func (r *Repository) GetByID(ctx context.Context, sessionID string) (*Session, error) {
	row := r.pool.QueryRow(ctx, `
SELECT
    id::text,
    domain_id,
    mode,
    session_goal,
    locale_language_id,
    locale_region_id,
    result_confidence,
    progress_state,
    started_at,
    finished_at,
    created_at,
    updated_at
FROM runtime.sessions
WHERE id = $1::uuid
`, sessionID)

	session, err := scanSession(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrSessionNotFound
	}

	return session, err
}

func (r *Repository) UpdateProgress(ctx context.Context, sessionID, progressState string, resultConfidence *float64, finishedAt *time.Time) error {
	commandTag, err := r.pool.Exec(ctx, `
UPDATE runtime.sessions
SET
    progress_state = $2,
    result_confidence = $3,
    finished_at = $4,
    updated_at = NOW()
WHERE id = $1::uuid
`, sessionID, progressState, resultConfidence, finishedAt)
	if err != nil {
		return fmt.Errorf("update session progress: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return ErrSessionNotFound
	}

	return nil
}

func scanSession(row pgx.Row) (*Session, error) {
	var session Session
	var goal sql.NullString
	var region sql.NullInt64
	var confidence sql.NullFloat64
	var finished sql.NullTime

	if err := row.Scan(
		&session.ID,
		&session.DomainID,
		&session.Mode,
		&goal,
		&session.LocaleLanguageID,
		&region,
		&confidence,
		&session.ProgressState,
		&session.StartedAt,
		&finished,
		&session.CreatedAt,
		&session.UpdatedAt,
	); err != nil {
		return nil, err
	}

	if goal.Valid {
		session.SessionGoal = goal.String
	}
	if region.Valid {
		value := region.Int64
		session.LocaleRegionID = &value
	}
	if confidence.Valid {
		value := confidence.Float64
		session.ResultConfidence = &value
	}
	if finished.Valid {
		value := finished.Time
		session.FinishedAt = &value
	}

	return &session, nil
}

func nullIfEmpty(value string) any {
	if value == "" {
		return nil
	}
	return value
}
