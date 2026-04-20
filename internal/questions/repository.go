package questions

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/andygellermen/CEE4AI/internal/packages"
)

var ErrQuestionNotFound = errors.New("question not found")

type QuestionMeta struct {
	ID           int64
	DomainID     int64
	QuestionType string
	ScoringMode  string
}

type DeliveryQuestion struct {
	ID                              int64
	ExternalID                      string
	QuestionFamily                  string
	QuestionType                    string
	IntendedUse                     string
	ConfidenceTier                  string
	EstimatedTimeSeconds            *int
	CognitiveLoadLevel              string
	MeaningDepth                    string
	WorldviewSensitivity            string
	SymbolicInterpretationRelevance string
	ExistentialLoadLevel            string
	IsSensitive                     bool
	AgeGate                         int
	Title                           string
	QuestionText                    string
	ExplanationText                 string
	RequiresHumanReview             bool
	WorldviewSensitive              bool
	Options                         []QuestionOption
}

type QuestionOption struct {
	ID           int64
	OptionKey    string
	OptionText   string
	DisplayOrder int
}

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) CountActiveByDomain(ctx context.Context, domainID int64) (int, error) {
	var count int
	if err := r.pool.QueryRow(ctx, `
SELECT COUNT(*)
FROM content.question_master
WHERE domain_id = $1
  AND is_active = TRUE
  AND review_status = 'active'
`, domainID).Scan(&count); err != nil {
		return 0, fmt.Errorf("count active questions: %w", err)
	}

	return count, nil
}

func (r *Repository) BuildPackage(ctx context.Context, domainID int64, offset, limit int) (*packages.QuestionPlan, error) {
	rows, err := r.pool.Query(ctx, `
SELECT id, COALESCE(estimated_time_seconds, 0)
FROM content.question_master
WHERE domain_id = $1
  AND is_active = TRUE
  AND review_status = 'active'
ORDER BY id
LIMIT $2 OFFSET $3
`, domainID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query package questions: %w", err)
	}
	defer rows.Close()

	plan := &packages.QuestionPlan{}
	for rows.Next() {
		var questionID int64
		var estimate int
		if err := rows.Scan(&questionID, &estimate); err != nil {
			return nil, fmt.Errorf("scan package question: %w", err)
		}
		plan.QuestionIDs = append(plan.QuestionIDs, questionID)
		plan.EstimatedTimeSeconds += estimate
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate package questions: %w", err)
	}

	return plan, nil
}

func (r *Repository) GetNextUnansweredQuestionID(ctx context.Context, sessionID string, domainID int64) (int64, error) {
	var questionID int64
	err := r.pool.QueryRow(ctx, `
SELECT qm.id
FROM content.question_master qm
WHERE qm.domain_id = $2
  AND qm.is_active = TRUE
  AND qm.review_status = 'active'
  AND NOT EXISTS (
      SELECT 1
      FROM runtime.answers a
      WHERE a.session_id = $1::uuid
        AND a.question_id = qm.id
  )
ORDER BY qm.id
LIMIT 1
`, sessionID, domainID).Scan(&questionID)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ErrQuestionNotFound
	}
	if err != nil {
		return 0, fmt.Errorf("query next unanswered question: %w", err)
	}

	return questionID, nil
}

func (r *Repository) GetQuestionPosition(ctx context.Context, domainID, questionID int64) (int, error) {
	var position int
	err := r.pool.QueryRow(ctx, `
WITH ranked AS (
    SELECT
        id,
        ROW_NUMBER() OVER (ORDER BY id) AS position
    FROM content.question_master
    WHERE domain_id = $1
      AND is_active = TRUE
      AND review_status = 'active'
)
SELECT position
FROM ranked
WHERE id = $2
`, domainID, questionID).Scan(&position)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ErrQuestionNotFound
	}
	if err != nil {
		return 0, fmt.Errorf("query question position: %w", err)
	}

	return position, nil
}

func (r *Repository) GetMetaByID(ctx context.Context, questionID int64) (*QuestionMeta, error) {
	var meta QuestionMeta
	err := r.pool.QueryRow(ctx, `
SELECT id, domain_id, question_type, scoring_mode
FROM content.question_master
WHERE id = $1
  AND is_active = TRUE
  AND review_status = 'active'
`, questionID).Scan(&meta.ID, &meta.DomainID, &meta.QuestionType, &meta.ScoringMode)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrQuestionNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("query question meta: %w", err)
	}

	return &meta, nil
}

func (r *Repository) GetByIDForLocale(ctx context.Context, questionID, languageID int64, regionID *int64) (*DeliveryQuestion, error) {
	var question DeliveryQuestion
	var title sql.NullString
	var explanation sql.NullString
	var estimated sql.NullInt64
	err := r.pool.QueryRow(ctx, `
SELECT
    qm.id,
    qm.external_id,
    qm.question_family,
    qm.question_type,
    COALESCE(qm.intended_use, ''),
    COALESCE(qm.confidence_tier, ''),
    qm.estimated_time_seconds,
    COALESCE(qm.cognitive_load_level, ''),
    COALESCE(qm.meaning_depth, ''),
    COALESCE(qm.worldview_sensitivity, ''),
    COALESCE(qm.symbolic_interpretation_relevance, ''),
    COALESCE(qm.existential_load_level, ''),
    qm.is_sensitive,
    qm.age_gate,
    qt.title,
    qt.question_text,
    qt.explanation_text,
    qt.requires_human_review,
    qt.worldview_sensitive
FROM content.question_master qm
JOIN LATERAL (
    SELECT
        title,
        question_text,
        explanation_text,
        requires_human_review,
        worldview_sensitive
    FROM content.question_translation
    WHERE question_id = qm.id
      AND language_id = $2
      AND is_active = TRUE
      AND (region_id IS NULL OR region_id IS NOT DISTINCT FROM $3)
    ORDER BY
      CASE WHEN region_id IS NOT DISTINCT FROM $3 THEN 0 ELSE 1 END,
      version DESC
    LIMIT 1
) qt ON TRUE
WHERE qm.id = $1
  AND qm.is_active = TRUE
  AND qm.review_status = 'active'
`, questionID, languageID, regionID).Scan(
		&question.ID,
		&question.ExternalID,
		&question.QuestionFamily,
		&question.QuestionType,
		&question.IntendedUse,
		&question.ConfidenceTier,
		&estimated,
		&question.CognitiveLoadLevel,
		&question.MeaningDepth,
		&question.WorldviewSensitivity,
		&question.SymbolicInterpretationRelevance,
		&question.ExistentialLoadLevel,
		&question.IsSensitive,
		&question.AgeGate,
		&title,
		&question.QuestionText,
		&explanation,
		&question.RequiresHumanReview,
		&question.WorldviewSensitive,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrQuestionNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("query localized question: %w", err)
	}

	if title.Valid {
		question.Title = title.String
	}
	if explanation.Valid {
		question.ExplanationText = explanation.String
	}
	if estimated.Valid {
		value := int(estimated.Int64)
		question.EstimatedTimeSeconds = &value
	}

	options, err := r.loadOptions(ctx, questionID, languageID, regionID)
	if err != nil {
		return nil, err
	}
	question.Options = options

	return &question, nil
}

func (r *Repository) loadOptions(ctx context.Context, questionID, languageID int64, regionID *int64) ([]QuestionOption, error) {
	rows, err := r.pool.Query(ctx, `
SELECT
    qom.id,
    qom.option_key,
    qot.option_text,
    qom.display_order
FROM content.question_option_master qom
JOIN LATERAL (
    SELECT option_text
    FROM content.question_option_translation
    WHERE option_id = qom.id
      AND language_id = $2
      AND is_active = TRUE
      AND (region_id IS NULL OR region_id IS NOT DISTINCT FROM $3)
    ORDER BY CASE WHEN region_id IS NOT DISTINCT FROM $3 THEN 0 ELSE 1 END
    LIMIT 1
) qot ON TRUE
WHERE qom.question_id = $1
  AND qom.is_active = TRUE
ORDER BY qom.display_order
`, questionID, languageID, regionID)
	if err != nil {
		return nil, fmt.Errorf("query question options: %w", err)
	}
	defer rows.Close()

	options := make([]QuestionOption, 0)
	for rows.Next() {
		var option QuestionOption
		if err := rows.Scan(&option.ID, &option.OptionKey, &option.OptionText, &option.DisplayOrder); err != nil {
			return nil, fmt.Errorf("scan question option: %w", err)
		}
		options = append(options, option)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate question options: %w", err)
	}

	return options, nil
}
