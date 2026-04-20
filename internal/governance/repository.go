package governance

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrRulesetNotFound = errors.New("ruleset version not found")

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) GetActiveRuleset(ctx context.Context, domainID int64, slug string) (*RulesetVersion, error) {
	row := r.pool.QueryRow(ctx, `
SELECT
    id,
    domain_id,
    slug,
    version,
    description,
    effective_from,
    effective_to,
    is_active
FROM governance.ruleset_versions
WHERE slug = $1
  AND is_active = TRUE
  AND effective_from <= NOW()
  AND (effective_to IS NULL OR effective_to > NOW())
  AND (domain_id = $2 OR domain_id IS NULL)
ORDER BY
    CASE WHEN domain_id = $2 THEN 0 ELSE 1 END,
    effective_from DESC,
    id DESC
LIMIT 1
`, slug, domainID)

	ruleset, err := scanRuleset(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrRulesetNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("query active ruleset: %w", err)
	}

	return ruleset, nil
}

func (r *Repository) ListActiveSensitivityRules(ctx context.Context) ([]SensitivityRule, error) {
	rows, err := r.pool.Query(ctx, `
SELECT id, slug, name, description, policy_payload, is_active
FROM governance.sensitivity_rules
WHERE is_active = TRUE
ORDER BY id
`)
	if err != nil {
		return nil, fmt.Errorf("query sensitivity rules: %w", err)
	}
	defer rows.Close()

	rules := make([]SensitivityRule, 0)
	for rows.Next() {
		var rule SensitivityRule
		var body []byte
		var description sql.NullString
		if err := rows.Scan(&rule.ID, &rule.Slug, &rule.Name, &description, &body, &rule.IsActive); err != nil {
			return nil, fmt.Errorf("scan sensitivity rule: %w", err)
		}
		if description.Valid {
			rule.Description = description.String
		}
		if len(body) > 0 {
			if err := json.Unmarshal(body, &rule.PolicyPayload); err != nil {
				return nil, fmt.Errorf("unmarshal sensitivity rule payload: %w", err)
			}
		}
		rules = append(rules, rule)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate sensitivity rules: %w", err)
	}

	return rules, nil
}

func (r *Repository) ListActiveContentPolicies(ctx context.Context, domainID int64, regionID *int64) ([]ContentPolicy, error) {
	rows, err := r.pool.Query(ctx, `
SELECT id, slug, name, description
FROM governance.content_policies
WHERE is_active = TRUE
  AND (applies_to_domain_id = $1 OR applies_to_domain_id IS NULL)
  AND (applies_to_region_id IS NULL OR applies_to_region_id IS NOT DISTINCT FROM $2)
ORDER BY
    CASE WHEN applies_to_domain_id = $1 THEN 0 ELSE 1 END,
    CASE WHEN applies_to_region_id IS NOT DISTINCT FROM $2 THEN 0 ELSE 1 END,
    id
`, domainID, regionID)
	if err != nil {
		return nil, fmt.Errorf("query content policies: %w", err)
	}
	defer rows.Close()

	policies := make([]ContentPolicy, 0)
	for rows.Next() {
		var policy ContentPolicy
		var description sql.NullString
		if err := rows.Scan(&policy.ID, &policy.Slug, &policy.Name, &description); err != nil {
			return nil, fmt.Errorf("scan content policy: %w", err)
		}
		if description.Valid {
			policy.Description = description.String
		}
		policies = append(policies, policy)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate content policies: %w", err)
	}

	return policies, nil
}

func (r *Repository) CreateAuditLog(ctx context.Context, params CreateAuditLogParams) error {
	body, err := json.Marshal(params.Payload)
	if err != nil {
		return fmt.Errorf("marshal audit payload: %w", err)
	}

	if _, err := r.pool.Exec(ctx, `
INSERT INTO governance.audit_log (actor_user_id, entity_type, entity_id, action, payload)
VALUES ($1::uuid, $2, $3, $4, $5)
`, params.ActorUserID, params.EntityType, params.EntityID, params.Action, body); err != nil {
		return fmt.Errorf("insert audit log: %w", err)
	}

	return nil
}

func (r *Repository) SummarizeSessionSensitivity(ctx context.Context, sessionID string) (*SessionSensitivitySummary, error) {
	var summary SessionSensitivitySummary
	if err := r.pool.QueryRow(ctx, `
SELECT
    COUNT(*) FILTER (WHERE qm.is_sensitive) AS answered_sensitive_questions,
    COUNT(*) FILTER (WHERE qm.age_gate > 0) AS answered_age_gated_questions,
    COUNT(*) FILTER (
        WHERE COALESCE(qt.worldview_sensitive, FALSE)
           OR LOWER(COALESCE(qm.worldview_sensitivity, '')) IN ('high', 'very_high')
    ) AS answered_worldview_sensitive_questions
FROM runtime.answers a
JOIN runtime.sessions s
  ON s.id = a.session_id
JOIN content.question_master qm
  ON qm.id = a.question_id
LEFT JOIN LATERAL (
    SELECT worldview_sensitive
    FROM content.question_translation
    WHERE question_id = qm.id
      AND language_id = s.locale_language_id
      AND is_active = TRUE
      AND localization_status = 'approved'
      AND (region_id IS NULL OR region_id IS NOT DISTINCT FROM s.locale_region_id)
    ORDER BY
      CASE WHEN region_id IS NOT DISTINCT FROM s.locale_region_id THEN 0 ELSE 1 END,
      version DESC
    LIMIT 1
) qt ON TRUE
WHERE a.session_id = $1::uuid
`, sessionID).Scan(
		&summary.AnsweredSensitiveQuestions,
		&summary.AnsweredAgeGatedQuestions,
		&summary.AnsweredWorldviewSensitiveQuestions,
	); err != nil {
		return nil, fmt.Errorf("summarize session sensitivity: %w", err)
	}

	return &summary, nil
}

func scanRuleset(row pgx.Row) (*RulesetVersion, error) {
	var ruleset RulesetVersion
	var domainID sql.NullInt64
	var description sql.NullString
	var effectiveTo sql.NullTime

	if err := row.Scan(
		&ruleset.ID,
		&domainID,
		&ruleset.Slug,
		&ruleset.Version,
		&description,
		&ruleset.EffectiveFrom,
		&effectiveTo,
		&ruleset.IsActive,
	); err != nil {
		return nil, err
	}

	if domainID.Valid {
		value := domainID.Int64
		ruleset.DomainID = &value
	}
	if description.Valid {
		ruleset.Description = description.String
	}
	if effectiveTo.Valid {
		value := effectiveTo.Time
		ruleset.EffectiveTo = &value
	}

	return &ruleset, nil
}
