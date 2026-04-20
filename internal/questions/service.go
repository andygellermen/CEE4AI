package questions

import (
	"context"
	"errors"

	"github.com/andygellermen/CEE4AI/internal/governance"
	"github.com/andygellermen/CEE4AI/internal/packages"
	"github.com/andygellermen/CEE4AI/internal/sessions"
)

type AnswerCounter interface {
	CountForSession(ctx context.Context, sessionID string) (int, error)
}

type NextQuestionResult struct {
	Session            *sessions.Session
	Package            *packages.SessionPackage
	Question           *DeliveryQuestion
	Governance         *governance.RuntimeDecision
	RemainingQuestions int
	HasMore            bool
}

type Service struct {
	sessionRepo       *sessions.Repository
	answerRepo        AnswerCounter
	questionRepo      *Repository
	packageService    *packages.Service
	governanceService *governance.Service
}

func NewService(
	sessionRepo *sessions.Repository,
	answerRepo AnswerCounter,
	questionRepo *Repository,
	packageService *packages.Service,
	governanceService *governance.Service,
) *Service {
	return &Service{
		sessionRepo:       sessionRepo,
		answerRepo:        answerRepo,
		questionRepo:      questionRepo,
		packageService:    packageService,
		governanceService: governanceService,
	}
}

func (s *Service) Next(ctx context.Context, sessionID string) (*NextQuestionResult, error) {
	session, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	totalQuestions, err := s.questionRepo.CountActiveByDomain(ctx, session.DomainID, session.LocaleLanguageID, session.LocaleRegionID)
	if err != nil {
		return nil, err
	}

	answeredQuestions, err := s.answerRepo.CountForSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	remaining := totalQuestions - answeredQuestions
	if remaining <= 0 {
		return &NextQuestionResult{
			Session:            session,
			RemainingQuestions: 0,
			HasMore:            false,
		}, nil
	}

	questionID, err := s.questionRepo.GetNextUnansweredQuestionID(ctx, sessionID, session.DomainID, session.LocaleLanguageID, session.LocaleRegionID)
	if err != nil {
		if errors.Is(err, ErrQuestionNotFound) {
			return &NextQuestionResult{
				Session:            session,
				RemainingQuestions: 0,
				HasMore:            false,
			}, nil
		}
		return nil, err
	}

	position, err := s.questionRepo.GetQuestionPosition(ctx, session.DomainID, session.LocaleLanguageID, questionID, session.LocaleRegionID)
	if err != nil {
		return nil, err
	}

	pkg, err := s.packageService.EnsureByQuestionPosition(ctx, session.ID, session.DomainID, session.LocaleLanguageID, session.LocaleRegionID, position)
	if err != nil {
		return nil, err
	}

	question, err := s.questionRepo.GetByIDForLocale(ctx, questionID, session.LocaleLanguageID, session.LocaleRegionID)
	if err != nil {
		return nil, err
	}

	decision, err := s.governanceService.ResolveQuestionDecision(ctx, governance.QuestionPolicyInput{
		QuestionID:           question.ID,
		DomainID:             session.DomainID,
		Mode:                 session.Mode,
		LocaleRegionID:       session.LocaleRegionID,
		QuestionFamily:       question.QuestionFamily,
		IsSensitive:          question.IsSensitive,
		AgeGate:              question.AgeGate,
		MeaningDepth:         question.MeaningDepth,
		WorldviewSensitivity: question.WorldviewSensitivity,
		RequiresHumanReview:  question.RequiresHumanReview,
		WorldviewSensitive:   question.WorldviewSensitive,
	})
	if err != nil {
		return nil, err
	}

	if err := s.governanceService.Audit(ctx, governance.CreateAuditLogParams{
		EntityType: "question",
		EntityID:   question.ExternalID,
		Action:     "runtime.question_presented",
		Payload: map[string]any{
			"session_id":        session.ID,
			"question_id":       question.ID,
			"question_family":   question.QuestionFamily,
			"delivery_mode":     decision.DeliveryMode,
			"severity":          decision.Severity,
			"sensitivity_flags": decision.SensitivityFlags,
			"ruleset_version":   decision.RulesetVersion,
			"review_required":   decision.ReviewRequired,
		},
	}); err != nil {
		return nil, err
	}

	return &NextQuestionResult{
		Session:            session,
		Package:            pkg,
		Question:           question,
		Governance:         decision,
		RemainingQuestions: remaining,
		HasMore:            true,
	}, nil
}
