package results

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) ReplaceProfileVector(ctx context.Context, sessionID, vectorType string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal profile vector %s: %w", vectorType, err)
	}

	if _, err := r.pool.Exec(ctx, `
DELETE FROM runtime.profile_vectors
WHERE session_id = $1::uuid
  AND vector_type = $2
`, sessionID, vectorType); err != nil {
		return fmt.Errorf("delete profile vector %s: %w", vectorType, err)
	}

	if _, err := r.pool.Exec(ctx, `
INSERT INTO runtime.profile_vectors (session_id, vector_type, vector_payload)
VALUES ($1::uuid, $2, $3)
`, sessionID, vectorType, body); err != nil {
		return fmt.Errorf("insert profile vector %s: %w", vectorType, err)
	}

	return nil
}

func (r *Repository) ReplaceSnapshot(
	ctx context.Context,
	sessionID, resultType, profileDepth, certaintyLevel, rulesetVersion string,
	payload any,
) (*Snapshot, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal snapshot payload: %w", err)
	}

	if _, err := r.pool.Exec(ctx, `
DELETE FROM runtime.result_snapshots
WHERE session_id = $1::uuid
  AND result_type = $2
`, sessionID, resultType); err != nil {
		return nil, fmt.Errorf("delete previous snapshot: %w", err)
	}

	row := r.pool.QueryRow(ctx, `
INSERT INTO runtime.result_snapshots (
    session_id,
    result_type,
    profile_depth,
    certainty_level,
    snapshot_payload,
    ruleset_version
)
VALUES ($1::uuid, $2, $3, $4, $5, $6)
RETURNING
    id,
    session_id::text,
    result_type,
    profile_depth,
    certainty_level,
    snapshot_payload,
    ruleset_version,
    created_at
`, sessionID, resultType, profileDepth, certaintyLevel, body, rulesetVersion)

	return scanSnapshot(row)
}

func scanSnapshot(row pgx.Row) (*Snapshot, error) {
	var snapshot Snapshot
	if err := row.Scan(
		&snapshot.ID,
		&snapshot.SessionID,
		&snapshot.ResultType,
		&snapshot.ProfileDepth,
		&snapshot.CertaintyLevel,
		&snapshot.SnapshotPayload,
		&snapshot.RulesetVersion,
		&snapshot.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &snapshot, nil
}
