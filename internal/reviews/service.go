package reviews

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/andygellermen/CEE4AI/internal/governance"
)

var ErrInvalidReviewStatus = errors.New("invalid review status")

type FlagQuestionRequest struct {
	QuestionID   int64
	SessionID    *string
	FlagSlug     string
	ReviewerRole string
	Comment      string
	Severity     string
}

type ApplyQuestionDecisionRequest struct {
	QuestionID int64
	NewStatus  string
	Reason     string
}

type RecordLocalizationReviewRequest struct {
	TranslationID int64
	Status        string
	Comment       string
}

type RuntimeFlagRequest struct {
	QuestionID      int64
	SessionID       string
	Decision        *governance.RuntimeDecision
	FreeTextPresent bool
}

type Service struct {
	repo              *Repository
	governanceService *governance.Service
}

func NewService(repo *Repository, governanceService *governance.Service) *Service {
	return &Service{
		repo:              repo,
		governanceService: governanceService,
	}
}

func (s *Service) FlagQuestion(ctx context.Context, req FlagQuestionRequest) (*QuestionReview, error) {
	flagID, err := s.repo.GetFlagBySlug(ctx, strings.TrimSpace(req.FlagSlug))
	if err != nil {
		return nil, err
	}

	reviewerRole := strings.TrimSpace(req.ReviewerRole)
	if reviewerRole == "" {
		reviewerRole = "system"
	}

	review, err := s.repo.CreateQuestionReview(
		ctx,
		req.QuestionID,
		req.SessionID,
		reviewerRole,
		req.FlagSlug,
		flagID,
		strings.TrimSpace(req.Comment),
		strings.TrimSpace(req.Severity),
	)
	if err != nil {
		return nil, err
	}

	if err := s.governanceService.Audit(ctx, governance.CreateAuditLogParams{
		EntityType: "question",
		EntityID:   strconv.FormatInt(req.QuestionID, 10),
		Action:     "review.question_flagged",
		Payload: map[string]any{
			"session_id":    req.SessionID,
			"flag_slug":     req.FlagSlug,
			"reviewer_role": reviewerRole,
			"severity":      strings.TrimSpace(req.Severity),
			"comment":       strings.TrimSpace(req.Comment),
			"review_record": review.ID,
		},
	}); err != nil {
		return nil, err
	}

	return review, nil
}

func (s *Service) ApplyQuestionDecision(ctx context.Context, req ApplyQuestionDecisionRequest) (*ReviewDecision, error) {
	newStatus := normalizeQuestionStatus(req.NewStatus)
	if newStatus == "" {
		return nil, ErrInvalidReviewStatus
	}

	oldStatus, err := s.repo.GetQuestionStatus(ctx, req.QuestionID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.UpdateQuestionStatus(ctx, req.QuestionID, newStatus); err != nil {
		return nil, err
	}

	decision, err := s.repo.CreateReviewDecision(ctx, req.QuestionID, oldStatus, newStatus, strings.TrimSpace(req.Reason))
	if err != nil {
		return nil, err
	}

	if err := s.governanceService.Audit(ctx, governance.CreateAuditLogParams{
		EntityType: "question",
		EntityID:   strconv.FormatInt(req.QuestionID, 10),
		Action:     "review.question_status_changed",
		Payload: map[string]any{
			"old_status":  oldStatus,
			"new_status":  newStatus,
			"reason":      strings.TrimSpace(req.Reason),
			"decision_id": decision.ID,
		},
	}); err != nil {
		return nil, err
	}

	return decision, nil
}

func (s *Service) RecordLocalizationReview(ctx context.Context, req RecordLocalizationReviewRequest) (*LocalizationReview, error) {
	newStatus := normalizeLocalizationStatus(req.Status)
	if newStatus == "" {
		return nil, ErrInvalidReviewStatus
	}

	target, err := s.repo.GetTranslationTarget(ctx, req.TranslationID)
	if err != nil {
		return nil, err
	}

	if target.RequiresHumanReview && newStatus == "approved" && target.LocalizationStatus == "machine_seeded" {
		newStatus = "human_reviewed"
	}

	if err := s.repo.UpdateTranslationStatus(ctx, req.TranslationID, newStatus); err != nil {
		return nil, err
	}

	review, err := s.repo.CreateLocalizationReview(ctx, req.TranslationID, newStatus, strings.TrimSpace(req.Comment))
	if err != nil {
		return nil, err
	}

	if err := s.governanceService.Audit(ctx, governance.CreateAuditLogParams{
		EntityType: "translation",
		EntityID:   strconv.FormatInt(req.TranslationID, 10),
		Action:     "review.localization_status_changed",
		Payload: map[string]any{
			"question_id":            target.QuestionID,
			"old_status":             target.LocalizationStatus,
			"new_status":             newStatus,
			"comment":                strings.TrimSpace(req.Comment),
			"requires_human_review":  target.RequiresHumanReview,
			"worldview_sensitive":    target.WorldviewSensitive,
			"localization_review_id": review.ID,
		},
	}); err != nil {
		return nil, err
	}

	return review, nil
}

func (s *Service) FlagRuntimeSensitivity(ctx context.Context, req RuntimeFlagRequest) error {
	if req.Decision == nil || req.Decision.DeliveryMode != "guarded" {
		return nil
	}

	var flagSlugs []string
	for _, flag := range req.Decision.SensitivityFlags {
		switch flag {
		case "sensitive_content":
			flagSlugs = append(flagSlugs, "sensitive_content")
		case "age_gated":
			flagSlugs = append(flagSlugs, "age_gate")
		case "worldview_sensitive":
			flagSlugs = append(flagSlugs, "worldview_sensitive")
		case "human_review_required":
			flagSlugs = append(flagSlugs, "human_review_required")
		}
	}
	if req.FreeTextPresent {
		flagSlugs = append(flagSlugs, "sensitive_free_text")
	}

	seen := make(map[string]struct{})
	for _, flagSlug := range flagSlugs {
		if _, ok := seen[flagSlug]; ok {
			continue
		}
		seen[flagSlug] = struct{}{}

		sessionID := req.SessionID
		_, err := s.FlagQuestion(ctx, FlagQuestionRequest{
			QuestionID:   req.QuestionID,
			SessionID:    &sessionID,
			FlagSlug:     flagSlug,
			ReviewerRole: "system",
			Comment:      runtimeFlagComment(flagSlug, req.Decision),
			Severity:     req.Decision.Severity,
		})
		if err != nil && !errors.Is(err, ErrReviewFlagNotFound) {
			return fmt.Errorf("flag runtime sensitivity %s: %w", flagSlug, err)
		}
	}

	return nil
}

func normalizeQuestionStatus(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "draft":
		return "draft"
	case "active":
		return "active"
	case "review":
		return "review"
	case "flagged":
		return "flagged"
	case "disabled":
		return "disabled"
	case "archived":
		return "archived"
	default:
		return ""
	}
}

func normalizeLocalizationStatus(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "draft":
		return "draft"
	case "machine_seeded":
		return "machine_seeded"
	case "human_reviewed":
		return "human_reviewed"
	case "approved":
		return "approved"
	case "archived":
		return "archived"
	default:
		return ""
	}
}

func runtimeFlagComment(flagSlug string, decision *governance.RuntimeDecision) string {
	switch flagSlug {
	case "age_gate":
		return "runtime flow stored an answer for age-gated content in guarded mode"
	case "worldview_sensitive":
		return "runtime flow stored worldview-sensitive content and attached guarded handling"
	case "human_review_required":
		return "runtime flow encountered content marked for human-reviewed handling"
	case "sensitive_free_text":
		return "runtime flow stored free-text input under guarded review conditions"
	default:
		return "runtime flow stored sensitive content under guarded handling"
	}
}
