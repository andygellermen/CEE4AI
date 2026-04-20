package reviews

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrReviewFlagNotFound  = errors.New("review flag not found")
	ErrQuestionNotFound    = errors.New("review question not found")
	ErrTranslationNotFound = errors.New("review translation not found")
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) GetFlagBySlug(ctx context.Context, slug string) (int64, error) {
	var id int64
	err := r.pool.QueryRow(ctx, `
SELECT id
FROM review.review_flags
WHERE slug = $1
  AND is_active = TRUE
`, slug).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ErrReviewFlagNotFound
	}
	if err != nil {
		return 0, fmt.Errorf("query review flag: %w", err)
	}

	return id, nil
}

func (r *Repository) CreateQuestionReview(ctx context.Context, questionID int64, sessionID *string, reviewerRole, flagSlug string, flagID int64, comment, severity string) (*QuestionReview, error) {
	row := r.pool.QueryRow(ctx, `
INSERT INTO review.question_reviews (
    question_id,
    session_id,
    reviewer_role,
    flag_id,
    comment,
    severity
)
VALUES ($1, $2::uuid, $3, $4, $5, $6)
RETURNING id, question_id, session_id::text, reviewer_role, flag_id, comment, severity, created_at
`, questionID, sessionID, nullIfEmpty(reviewerRole), flagID, nullIfEmpty(comment), nullIfEmpty(severity))

	return scanQuestionReview(row, flagSlug)
}

func (r *Repository) GetQuestionStatus(ctx context.Context, questionID int64) (string, error) {
	var status string
	err := r.pool.QueryRow(ctx, `
SELECT review_status
FROM content.question_master
WHERE id = $1
`, questionID).Scan(&status)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrQuestionNotFound
	}
	if err != nil {
		return "", fmt.Errorf("query question review status: %w", err)
	}

	return status, nil
}

func (r *Repository) UpdateQuestionStatus(ctx context.Context, questionID int64, status string) error {
	tag, err := r.pool.Exec(ctx, `
UPDATE content.question_master
SET review_status = $2,
    updated_at = NOW()
WHERE id = $1
`, questionID, status)
	if err != nil {
		return fmt.Errorf("update question review status: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrQuestionNotFound
	}
	return nil
}

func (r *Repository) CreateReviewDecision(ctx context.Context, questionID int64, oldStatus, newStatus, reason string) (*ReviewDecision, error) {
	row := r.pool.QueryRow(ctx, `
INSERT INTO review.review_decisions (
    question_id,
    old_status,
    new_status,
    reason
)
VALUES ($1, $2, $3, $4)
RETURNING id, question_id, old_status, new_status, reason, created_at
`, questionID, nullIfEmpty(oldStatus), nullIfEmpty(newStatus), nullIfEmpty(reason))

	var decision ReviewDecision
	var oldValue sql.NullString
	var reasonValue sql.NullString
	if err := row.Scan(&decision.ID, &decision.QuestionID, &oldValue, &decision.NewStatus, &reasonValue, &decision.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrQuestionNotFound
		}
		return nil, fmt.Errorf("insert review decision: %w", err)
	}
	if oldValue.Valid {
		decision.OldStatus = oldValue.String
	}
	if reasonValue.Valid {
		decision.Reason = reasonValue.String
	}
	return &decision, nil
}

func (r *Repository) GetTranslationTarget(ctx context.Context, translationID int64) (*TranslationTarget, error) {
	row := r.pool.QueryRow(ctx, `
SELECT id, question_id, localization_status, requires_human_review, worldview_sensitive
FROM content.question_translation
WHERE id = $1
`, translationID)

	var target TranslationTarget
	if err := row.Scan(
		&target.ID,
		&target.QuestionID,
		&target.LocalizationStatus,
		&target.RequiresHumanReview,
		&target.WorldviewSensitive,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTranslationNotFound
		}
		return nil, fmt.Errorf("query translation target: %w", err)
	}

	return &target, nil
}

func (r *Repository) UpdateTranslationStatus(ctx context.Context, translationID int64, status string) error {
	tag, err := r.pool.Exec(ctx, `
UPDATE content.question_translation
SET localization_status = $2,
    updated_at = NOW()
WHERE id = $1
`, translationID, status)
	if err != nil {
		return fmt.Errorf("update localization status: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrTranslationNotFound
	}
	return nil
}

func (r *Repository) CreateLocalizationReview(ctx context.Context, translationID int64, status, comment string) (*LocalizationReview, error) {
	row := r.pool.QueryRow(ctx, `
INSERT INTO review.localization_reviews (
    translation_id,
    status,
    comment
)
VALUES ($1, $2, $3)
RETURNING id, translation_id, status, comment, created_at
`, translationID, status, nullIfEmpty(comment))

	var review LocalizationReview
	var commentValue sql.NullString
	if err := row.Scan(&review.ID, &review.TranslationID, &review.Status, &commentValue, &review.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTranslationNotFound
		}
		return nil, fmt.Errorf("insert localization review: %w", err)
	}
	if commentValue.Valid {
		review.Comment = commentValue.String
	}
	return &review, nil
}

func (r *Repository) CountForSession(ctx context.Context, sessionID string) (int, error) {
	var count int
	if err := r.pool.QueryRow(ctx, `
SELECT COUNT(*)
FROM review.question_reviews
WHERE session_id = $1::uuid
`, sessionID).Scan(&count); err != nil {
		return 0, fmt.Errorf("count question reviews for session: %w", err)
	}
	return count, nil
}

func scanQuestionReview(row pgx.Row, flagSlug string) (*QuestionReview, error) {
	var review QuestionReview
	var sessionID sql.NullString
	var reviewerRole sql.NullString
	var comment sql.NullString
	var severity sql.NullString
	if err := row.Scan(
		&review.ID,
		&review.QuestionID,
		&sessionID,
		&reviewerRole,
		&review.FlagID,
		&comment,
		&severity,
		&review.CreatedAt,
	); err != nil {
		return nil, fmt.Errorf("scan question review: %w", err)
	}

	review.FlagSlug = flagSlug
	if sessionID.Valid {
		value := sessionID.String
		review.SessionID = &value
	}
	if reviewerRole.Valid {
		review.ReviewerRole = reviewerRole.String
	}
	if comment.Valid {
		review.Comment = comment.String
	}
	if severity.Valid {
		review.Severity = severity.String
	}

	return &review, nil
}

func nullIfEmpty(value string) any {
	if value == "" {
		return nil
	}
	return value
}
