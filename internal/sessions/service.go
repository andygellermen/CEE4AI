package sessions

import (
	"context"
	"errors"
	"strings"

	"github.com/andygellermen/CEE4AI/internal/packages"
)

const (
	DefaultMode                   = "snapshot"
	DefaultLocaleLanguageID int64 = 1
	DefaultProgressState          = "active"
)

var ErrInvalidMode = errors.New("invalid session mode")

type StartSessionRequest struct {
	DomainID         int64
	Mode             string
	SessionGoal      string
	LocaleLanguageID *int64
	LocaleRegionID   *int64
}

type StartSessionResult struct {
	Session      *Session
	FirstPackage *packages.SessionPackage
}

type Service struct {
	repo           *Repository
	packageService *packages.Service
}

func NewService(repo *Repository, packageService *packages.Service) *Service {
	return &Service{
		repo:           repo,
		packageService: packageService,
	}
}

func (s *Service) Start(ctx context.Context, req StartSessionRequest) (*StartSessionResult, error) {
	mode := strings.TrimSpace(req.Mode)
	if mode == "" {
		mode = DefaultMode
	}

	if mode != "snapshot" && mode != "guided_progression" {
		return nil, ErrInvalidMode
	}

	languageID := DefaultLocaleLanguageID
	if req.LocaleLanguageID != nil {
		languageID = *req.LocaleLanguageID
	}

	session, err := s.repo.Create(ctx, CreateSessionParams{
		DomainID:         req.DomainID,
		Mode:             mode,
		SessionGoal:      strings.TrimSpace(req.SessionGoal),
		LocaleLanguageID: languageID,
		LocaleRegionID:   req.LocaleRegionID,
		ProgressState:    DefaultProgressState,
	})
	if err != nil {
		return nil, err
	}

	firstPackage, err := s.packageService.EnsureByQuestionPosition(ctx, session.ID, session.DomainID, 1)
	if err != nil {
		return nil, err
	}

	return &StartSessionResult{
		Session:      session,
		FirstPackage: firstPackage,
	}, nil
}

func (s *Service) Get(ctx context.Context, sessionID string) (*Session, error) {
	return s.repo.GetByID(ctx, sessionID)
}
