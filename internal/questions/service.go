package questions

import (
	"context"
	"errors"

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
	RemainingQuestions int
	HasMore            bool
}

type Service struct {
	sessionRepo    *sessions.Repository
	answerRepo     AnswerCounter
	questionRepo   *Repository
	packageService *packages.Service
}

func NewService(
	sessionRepo *sessions.Repository,
	answerRepo AnswerCounter,
	questionRepo *Repository,
	packageService *packages.Service,
) *Service {
	return &Service{
		sessionRepo:    sessionRepo,
		answerRepo:     answerRepo,
		questionRepo:   questionRepo,
		packageService: packageService,
	}
}

func (s *Service) Next(ctx context.Context, sessionID string) (*NextQuestionResult, error) {
	session, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	totalQuestions, err := s.questionRepo.CountActiveByDomain(ctx, session.DomainID)
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

	questionID, err := s.questionRepo.GetNextUnansweredQuestionID(ctx, sessionID, session.DomainID)
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

	position, err := s.questionRepo.GetQuestionPosition(ctx, session.DomainID, questionID)
	if err != nil {
		return nil, err
	}

	pkg, err := s.packageService.EnsureByQuestionPosition(ctx, session.ID, session.DomainID, position)
	if err != nil {
		return nil, err
	}

	question, err := s.questionRepo.GetByIDForLocale(ctx, questionID, session.LocaleLanguageID, session.LocaleRegionID)
	if err != nil {
		return nil, err
	}

	return &NextQuestionResult{
		Session:            session,
		Package:            pkg,
		Question:           question,
		RemainingQuestions: remaining,
		HasMore:            true,
	}, nil
}
