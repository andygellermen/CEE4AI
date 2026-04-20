package governance

import (
	"context"
	"fmt"
	"strings"
)

const (
	defaultRulesetVersion = "unversioned"
	deliveryModeStandard  = "standard"
	deliveryModeGuarded   = "guarded"
	severityLow           = "low"
	severityMedium        = "medium"
	severityHigh          = "high"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ResolveQuestionDecision(ctx context.Context, input QuestionPolicyInput) (*RuntimeDecision, error) {
	rulesetSlug := rulesetSlugForMode(input.Mode)
	rulesetVersion := defaultRulesetVersion

	ruleset, err := s.repo.GetActiveRuleset(ctx, input.DomainID, rulesetSlug)
	if err != nil && err != ErrRulesetNotFound {
		return nil, err
	}
	if ruleset != nil {
		rulesetVersion = ruleset.Version
	}

	decision := &RuntimeDecision{
		RulesetSlug:         rulesetSlug,
		RulesetVersion:      rulesetVersion,
		DeliveryMode:        deliveryModeStandard,
		Severity:            severityLow,
		ApprovedForDelivery: true,
	}

	policies, err := s.repo.ListActiveContentPolicies(ctx, input.DomainID, input.LocaleRegionID)
	if err != nil {
		return nil, err
	}
	for _, policy := range policies {
		decision.AppliedPolicies = append(decision.AppliedPolicies, policy.Slug)
	}

	rules, err := s.repo.ListActiveSensitivityRules(ctx)
	if err != nil {
		return nil, err
	}
	for _, rule := range rules {
		applyRule(decision, rule, input)
	}

	if input.IsSensitive && !contains(decision.SensitivityFlags, "sensitive_content") {
		decision.SensitivityFlags = append(decision.SensitivityFlags, "sensitive_content")
	}
	if input.AgeGate > 0 && !contains(decision.SensitivityFlags, "age_gated") {
		decision.SensitivityFlags = append(decision.SensitivityFlags, "age_gated")
	}
	if isWorldviewSensitive(input) && !contains(decision.SensitivityFlags, "worldview_sensitive") {
		decision.SensitivityFlags = append(decision.SensitivityFlags, "worldview_sensitive")
	}
	if input.RequiresHumanReview && !contains(decision.SensitivityFlags, "human_review_required") {
		decision.SensitivityFlags = append(decision.SensitivityFlags, "human_review_required")
	}

	if len(decision.SensitivityFlags) > 0 {
		decision.DeliveryMode = deliveryModeGuarded
		if decision.Severity == severityLow {
			decision.Severity = severityMedium
		}
	}

	if input.RequiresHumanReview {
		decision.ReviewRequired = true
	}
	if contains(decision.SensitivityFlags, "age_gated") {
		decision.Notes = append(decision.Notes, fmt.Sprintf("question carries an age gate of %d+", input.AgeGate))
	}
	if contains(decision.SensitivityFlags, "human_review_required") {
		decision.Notes = append(decision.Notes, "localized wording is marked for human-reviewed handling")
	}
	if contains(decision.SensitivityFlags, "worldview_sensitive") {
		decision.Notes = append(decision.Notes, "worldview-sensitive handling is enabled for this question")
	}

	return decision, nil
}

func (s *Service) SummarizeSession(ctx context.Context, sessionID string, domainID int64, mode string, reviewFlagCount int) (*SessionGovernanceSummary, error) {
	rulesetSlug := rulesetSlugForMode(mode)
	rulesetVersion := defaultRulesetVersion

	ruleset, err := s.repo.GetActiveRuleset(ctx, domainID, rulesetSlug)
	if err != nil && err != ErrRulesetNotFound {
		return nil, err
	}
	if ruleset != nil {
		rulesetVersion = ruleset.Version
	}

	sensitivity, err := s.repo.SummarizeSessionSensitivity(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	summary := &SessionGovernanceSummary{
		RulesetSlug:                         rulesetSlug,
		RulesetVersion:                      rulesetVersion,
		DeliveryMode:                        deliveryModeStandard,
		Severity:                            severityLow,
		AnsweredSensitiveQuestions:          sensitivity.AnsweredSensitiveQuestions,
		AnsweredAgeGatedQuestions:           sensitivity.AnsweredAgeGatedQuestions,
		AnsweredWorldviewSensitiveQuestions: sensitivity.AnsweredWorldviewSensitiveQuestions,
		ReviewFlagCount:                     reviewFlagCount,
	}

	if sensitivity.AnsweredSensitiveQuestions > 0 || sensitivity.AnsweredAgeGatedQuestions > 0 || sensitivity.AnsweredWorldviewSensitiveQuestions > 0 || reviewFlagCount > 0 {
		summary.DeliveryMode = deliveryModeGuarded
		summary.Severity = severityMedium
	}
	if sensitivity.AnsweredWorldviewSensitiveQuestions > 0 {
		summary.Severity = severityHigh
		summary.Guardrails = append(summary.Guardrails, "result language should remain non-dogmatic and reflective")
	}
	if sensitivity.AnsweredAgeGatedQuestions > 0 {
		summary.Guardrails = append(summary.Guardrails, "age-gated content was included and should remain carefully framed")
	}
	if reviewFlagCount > 0 {
		summary.Guardrails = append(summary.Guardrails, "session contains flagged review context for follow-up")
	}

	return summary, nil
}

func (s *Service) Audit(ctx context.Context, params CreateAuditLogParams) error {
	return s.repo.CreateAuditLog(ctx, params)
}

func rulesetSlugForMode(mode string) string {
	switch strings.TrimSpace(mode) {
	case "guided_progression":
		return "runtime.guided_progression"
	default:
		return "runtime.snapshot"
	}
}

func applyRule(decision *RuntimeDecision, rule SensitivityRule, input QuestionPolicyInput) {
	trigger := strings.TrimSpace(stringValue(rule.PolicyPayload["trigger"]))
	if trigger == "" {
		return
	}
	if !matchesTrigger(trigger, input) {
		return
	}

	flag := strings.TrimSpace(stringValue(rule.PolicyPayload["flag"]))
	if flag != "" && !contains(decision.SensitivityFlags, flag) {
		decision.SensitivityFlags = append(decision.SensitivityFlags, flag)
	}

	policyMode := strings.TrimSpace(stringValue(rule.PolicyPayload["delivery_mode"]))
	if policyMode == deliveryModeGuarded {
		decision.DeliveryMode = deliveryModeGuarded
	}

	severity := normalizeSeverity(stringValue(rule.PolicyPayload["severity"]))
	if severityRank(severity) > severityRank(decision.Severity) {
		decision.Severity = severity
	}

	if boolValue(rule.PolicyPayload["review_required"]) {
		decision.ReviewRequired = true
	}

	note := strings.TrimSpace(stringValue(rule.PolicyPayload["note"]))
	if note != "" && !contains(decision.Notes, note) {
		decision.Notes = append(decision.Notes, note)
	}
}

func matchesTrigger(trigger string, input QuestionPolicyInput) bool {
	switch trigger {
	case "is_sensitive":
		return input.IsSensitive
	case "age_gate":
		return input.AgeGate > 0
	case "requires_human_review":
		return input.RequiresHumanReview
	case "worldview_sensitive":
		return isWorldviewSensitive(input)
	default:
		return false
	}
}

func isWorldviewSensitive(input QuestionPolicyInput) bool {
	level := strings.ToLower(strings.TrimSpace(input.WorldviewSensitivity))
	return input.WorldviewSensitive || level == "high" || level == "very_high"
}

func normalizeSeverity(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case severityHigh:
		return severityHigh
	case severityMedium:
		return severityMedium
	default:
		return severityLow
	}
}

func severityRank(value string) int {
	switch normalizeSeverity(value) {
	case severityHigh:
		return 3
	case severityMedium:
		return 2
	default:
		return 1
	}
}

func stringValue(value any) string {
	text, _ := value.(string)
	return text
}

func boolValue(value any) bool {
	flag, _ := value.(bool)
	return flag
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
