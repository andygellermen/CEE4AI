package governance

import "time"

type RulesetVersion struct {
	ID            int64
	DomainID      *int64
	Slug          string
	Version       string
	Description   string
	EffectiveFrom time.Time
	EffectiveTo   *time.Time
	IsActive      bool
}

type SensitivityRule struct {
	ID            int64
	Slug          string
	Name          string
	Description   string
	PolicyPayload map[string]any
	IsActive      bool
}

type ContentPolicy struct {
	ID          int64
	Slug        string
	Name        string
	Description string
}

type QuestionPolicyInput struct {
	QuestionID           int64
	DomainID             int64
	Mode                 string
	LocaleRegionID       *int64
	QuestionFamily       string
	IsSensitive          bool
	AgeGate              int
	MeaningDepth         string
	WorldviewSensitivity string
	RequiresHumanReview  bool
	WorldviewSensitive   bool
}

type RuntimeDecision struct {
	RulesetSlug         string   `json:"ruleset_slug"`
	RulesetVersion      string   `json:"ruleset_version"`
	DeliveryMode        string   `json:"delivery_mode"`
	Severity            string   `json:"severity"`
	ReviewRequired      bool     `json:"review_required"`
	ApprovedForDelivery bool     `json:"approved_for_delivery"`
	SensitivityFlags    []string `json:"sensitivity_flags,omitempty"`
	AppliedPolicies     []string `json:"applied_policies,omitempty"`
	Notes               []string `json:"notes,omitempty"`
}

type SessionSensitivitySummary struct {
	AnsweredSensitiveQuestions          int
	AnsweredAgeGatedQuestions           int
	AnsweredWorldviewSensitiveQuestions int
}

type SessionGovernanceSummary struct {
	RulesetSlug                         string   `json:"ruleset_slug"`
	RulesetVersion                      string   `json:"ruleset_version"`
	DeliveryMode                        string   `json:"delivery_mode"`
	Severity                            string   `json:"severity"`
	AnsweredSensitiveQuestions          int      `json:"answered_sensitive_questions"`
	AnsweredAgeGatedQuestions           int      `json:"answered_age_gated_questions"`
	AnsweredWorldviewSensitiveQuestions int      `json:"answered_worldview_sensitive_questions"`
	ReviewFlagCount                     int      `json:"review_flag_count"`
	Guardrails                          []string `json:"guardrails,omitempty"`
}

type CreateAuditLogParams struct {
	ActorUserID *string
	EntityType  string
	EntityID    string
	Action      string
	Payload     any
}
