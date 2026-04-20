package results

import (
	"context"
	"strconv"
	"time"

	"github.com/andygellermen/CEE4AI/internal/answers"
	"github.com/andygellermen/CEE4AI/internal/governance"
	"github.com/andygellermen/CEE4AI/internal/packages"
	"github.com/andygellermen/CEE4AI/internal/questions"
	"github.com/andygellermen/CEE4AI/internal/reviews"
	"github.com/andygellermen/CEE4AI/internal/scoring"
	"github.com/andygellermen/CEE4AI/internal/sessions"
)

const (
	snapshotResultType   = "snapshot_profile"
	snapshotProfileDepth = "snapshot"
)

type Service struct {
	sessionRepo       *sessions.Repository
	answerRepo        *answers.Repository
	questionRepo      *questions.Repository
	resultRepo        *Repository
	scoringService    *scoring.Service
	governanceService *governance.Service
	reviewRepo        *reviews.Repository
}

func NewService(
	sessionRepo *sessions.Repository,
	answerRepo *answers.Repository,
	questionRepo *questions.Repository,
	resultRepo *Repository,
	scoringService *scoring.Service,
	governanceService *governance.Service,
	reviewRepo *reviews.Repository,
) *Service {
	return &Service{
		sessionRepo:       sessionRepo,
		answerRepo:        answerRepo,
		questionRepo:      questionRepo,
		resultRepo:        resultRepo,
		scoringService:    scoringService,
		governanceService: governanceService,
		reviewRepo:        reviewRepo,
	}
}

func (s *Service) BuildSnapshot(ctx context.Context, sessionID string) (*SnapshotResult, error) {
	session, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	totalQuestions, err := s.questionRepo.CountActiveByDomain(ctx, session.DomainID, session.LocaleLanguageID, session.LocaleRegionID)
	if err != nil {
		return nil, err
	}

	answeredQuestions, err := s.answerRepo.CountForSession(ctx, session.ID)
	if err != nil {
		return nil, err
	}

	vectors, err := s.scoringService.BuildSessionVectors(ctx, session.ID)
	if err != nil {
		return nil, err
	}

	confidence := calculateSnapshotConfidence(answeredQuestions, totalQuestions)
	certaintyLevel := classifyConfidence(confidence)
	completionRatio := calculateCompletionRatio(answeredQuestions, totalQuestions)
	progressState := sessions.DefaultProgressState
	var finishedAt *time.Time
	if totalQuestions > 0 && answeredQuestions >= totalQuestions {
		progressState = "completed"
		now := time.Now().UTC()
		finishedAt = &now
	}

	if err := s.sessionRepo.UpdateProgress(ctx, session.ID, progressState, &confidence, finishedAt); err != nil {
		return nil, err
	}

	if err := s.persistVectors(ctx, session.ID, vectors); err != nil {
		return nil, err
	}

	reviewFlagCount, err := s.reviewRepo.CountForSession(ctx, session.ID)
	if err != nil {
		return nil, err
	}

	governanceSummary, err := s.governanceService.SummarizeSession(ctx, session.ID, session.DomainID, session.Mode, reviewFlagCount)
	if err != nil {
		return nil, err
	}

	payload := &SnapshotPayload{
		SessionID:         session.ID,
		DomainID:          session.DomainID,
		Mode:              session.Mode,
		ProfileDepth:      snapshotProfileDepth,
		ProgressState:     progressState,
		AnsweredQuestions: answeredQuestions,
		TotalQuestions:    totalQuestions,
		CompletionRatio:   completionRatio,
		ResultConfidence:  confidence,
		CertaintyLevel:    certaintyLevel,
		Vectors:           vectors,
		Governance:        governanceSummary,
		TopSignals: map[string]string{
			"denktype":  scoring.TopSignal(vectors.Denktype),
			"skill":     scoring.TopSignal(vectors.Skill),
			"trait":     scoring.TopSignal(vectors.Trait),
			"meaning":   scoring.TopSignal(vectors.Meaning),
			"worldview": scoring.TopSignal(vectors.Worldview),
		},
	}

	snapshot, err := s.resultRepo.ReplaceSnapshot(
		ctx,
		session.ID,
		snapshotResultType,
		snapshotProfileDepth,
		certaintyLevel,
		governanceSummary.RulesetVersion,
		payload,
	)
	if err != nil {
		return nil, err
	}

	if err := s.governanceService.Audit(ctx, governance.CreateAuditLogParams{
		EntityType: "snapshot",
		EntityID:   strconv.FormatInt(snapshot.ID, 10),
		Action:     "runtime.snapshot_generated",
		Payload: map[string]any{
			"session_id":         session.ID,
			"result_type":        snapshot.ResultType,
			"ruleset_version":    governanceSummary.RulesetVersion,
			"review_flag_count":  governanceSummary.ReviewFlagCount,
			"delivery_mode":      governanceSummary.DeliveryMode,
			"answered_questions": answeredQuestions,
			"answered_sensitive": governanceSummary.AnsweredSensitiveQuestions,
			"answered_worldview": governanceSummary.AnsweredWorldviewSensitiveQuestions,
		},
	}); err != nil {
		return nil, err
	}

	return &SnapshotResult{
		Snapshot: snapshot,
		Payload:  payload,
	}, nil
}

func (s *Service) persistVectors(ctx context.Context, sessionID string, vectors *scoring.SnapshotVectors) error {
	if err := s.resultRepo.ReplaceProfileVector(ctx, sessionID, "denktype", vectors.Denktype); err != nil {
		return err
	}
	if err := s.resultRepo.ReplaceProfileVector(ctx, sessionID, "skill", vectors.Skill); err != nil {
		return err
	}
	if err := s.resultRepo.ReplaceProfileVector(ctx, sessionID, "trait", vectors.Trait); err != nil {
		return err
	}
	if err := s.resultRepo.ReplaceProfileVector(ctx, sessionID, "meaning", vectors.Meaning); err != nil {
		return err
	}
	if err := s.resultRepo.ReplaceProfileVector(ctx, sessionID, "worldview", vectors.Worldview); err != nil {
		return err
	}

	return nil
}

func calculateSnapshotConfidence(answeredQuestions, totalQuestions int) float64 {
	if answeredQuestions <= 0 || totalQuestions <= 0 {
		return 0
	}

	target := packages.DefaultPackageSize
	if totalQuestions < target {
		target = totalQuestions
	}

	confidence := float64(answeredQuestions) / float64(target)
	if confidence > 1 {
		return 1
	}
	return confidence
}

func calculateCompletionRatio(answeredQuestions, totalQuestions int) float64 {
	if answeredQuestions <= 0 || totalQuestions <= 0 {
		return 0
	}
	return float64(answeredQuestions) / float64(totalQuestions)
}

func classifyConfidence(confidence float64) string {
	switch {
	case confidence >= 0.75:
		return "high"
	case confidence >= 0.4:
		return "medium"
	default:
		return "low"
	}
}
