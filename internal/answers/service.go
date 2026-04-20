package answers

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/andygellermen/CEE4AI/internal/governance"
	"github.com/andygellermen/CEE4AI/internal/packages"
	"github.com/andygellermen/CEE4AI/internal/questions"
	"github.com/andygellermen/CEE4AI/internal/reviews"
	"github.com/andygellermen/CEE4AI/internal/scoring"
	"github.com/andygellermen/CEE4AI/internal/sessions"
)

var (
	ErrQuestionAlreadyAnswered = errors.New("question already answered")
	ErrQuestionSessionMismatch = errors.New("question does not belong to session domain")
)

type SubmitAnswerRequest struct {
	SessionID         string
	QuestionID        int64
	SelectedOptionIDs []int64
	ScaleValue        *int
	FreeTextAnswer    string
	CertaintyLevel    string
}

type SubmitAnswerResult struct {
	Answer            *Answer
	Package           *packages.SessionPackage
	Governance        *governance.RuntimeDecision
	ProgressState     string
	AnsweredQuestions int
	TotalQuestions    int
	HasMore           bool
}

type Service struct {
	sessionRepo       *sessions.Repository
	questionRepo      *questions.Repository
	answerRepo        *Repository
	packageService    *packages.Service
	scoringService    *scoring.Service
	governanceService *governance.Service
	reviewService     *reviews.Service
}

func NewService(
	sessionRepo *sessions.Repository,
	questionRepo *questions.Repository,
	answerRepo *Repository,
	packageService *packages.Service,
	scoringService *scoring.Service,
	governanceService *governance.Service,
	reviewService *reviews.Service,
) *Service {
	return &Service{
		sessionRepo:       sessionRepo,
		questionRepo:      questionRepo,
		answerRepo:        answerRepo,
		packageService:    packageService,
		scoringService:    scoringService,
		governanceService: governanceService,
		reviewService:     reviewService,
	}
}

func (s *Service) Submit(ctx context.Context, req SubmitAnswerRequest) (*SubmitAnswerResult, error) {
	session, err := s.sessionRepo.GetByID(ctx, req.SessionID)
	if err != nil {
		return nil, err
	}

	meta, err := s.questionRepo.GetMetaByID(ctx, req.QuestionID)
	if err != nil {
		return nil, err
	}

	if meta.DomainID != session.DomainID {
		return nil, ErrQuestionSessionMismatch
	}

	question, err := s.questionRepo.GetByIDForLocale(ctx, req.QuestionID, session.LocaleLanguageID, session.LocaleRegionID)
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

	exists, err := s.answerRepo.ExistsForSessionQuestion(ctx, req.SessionID, req.QuestionID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrQuestionAlreadyAnswered
	}

	position, err := s.questionRepo.GetQuestionPosition(ctx, session.DomainID, session.LocaleLanguageID, req.QuestionID, session.LocaleRegionID)
	if err != nil {
		return nil, err
	}

	pkg, err := s.packageService.EnsureByQuestionPosition(ctx, session.ID, session.DomainID, session.LocaleLanguageID, session.LocaleRegionID, position)
	if err != nil {
		return nil, err
	}

	evaluation, err := s.scoringService.EvaluateAnswer(ctx, scoring.EvaluateAnswerParams{
		QuestionID:        req.QuestionID,
		QuestionType:      meta.QuestionType,
		SelectedOptionIDs: req.SelectedOptionIDs,
		ScaleValue:        req.ScaleValue,
		FreeTextAnswer:    strings.TrimSpace(req.FreeTextAnswer),
	})
	if err != nil {
		return nil, err
	}

	selectedOptionIDs, err := scoring.MarshalSelectedOptionIDs(evaluation.SelectedOptionIDs)
	if err != nil {
		return nil, err
	}

	var packageID *int64
	if pkg != nil {
		packageID = &pkg.ID
	}

	answer, err := s.answerRepo.Create(ctx, CreateAnswerParams{
		SessionID:         session.ID,
		PackageID:         packageID,
		QuestionID:        req.QuestionID,
		AnswerKind:        evaluation.AnswerKind,
		SelectedOptionIDs: selectedOptionIDs,
		ScaleValue:        evaluation.ScaleValue,
		FreeTextAnswer:    evaluation.FreeTextAnswer,
		RawScore:          evaluation.RawScore,
		EvaluatedScore:    evaluation.EvaluatedScore,
		CertaintyLevel:    strings.TrimSpace(req.CertaintyLevel),
	})
	if err != nil {
		return nil, err
	}

	answeredQuestions, err := s.answerRepo.CountForSession(ctx, session.ID)
	if err != nil {
		return nil, err
	}

	totalQuestions, err := s.questionRepo.CountActiveByDomain(ctx, session.DomainID, session.LocaleLanguageID, session.LocaleRegionID)
	if err != nil {
		return nil, err
	}

	progressState := sessions.DefaultProgressState
	var finishedAt *time.Time
	if answeredQuestions >= totalQuestions && totalQuestions > 0 {
		progressState = "completed"
		now := time.Now().UTC()
		finishedAt = &now
	}

	if err := s.sessionRepo.UpdateProgress(ctx, session.ID, progressState, session.ResultConfidence, finishedAt); err != nil {
		return nil, err
	}

	if err := s.governanceService.Audit(ctx, governance.CreateAuditLogParams{
		EntityType: "answer",
		EntityID:   strconv.FormatInt(answer.ID, 10),
		Action:     "runtime.answer_submitted",
		Payload: map[string]any{
			"session_id":        session.ID,
			"question_id":       req.QuestionID,
			"answer_kind":       answer.AnswerKind,
			"delivery_mode":     decision.DeliveryMode,
			"severity":          decision.Severity,
			"sensitivity_flags": decision.SensitivityFlags,
			"review_required":   decision.ReviewRequired,
			"progress_state":    progressState,
		},
	}); err != nil {
		return nil, err
	}

	if err := s.reviewService.FlagRuntimeSensitivity(ctx, reviews.RuntimeFlagRequest{
		QuestionID:      req.QuestionID,
		SessionID:       session.ID,
		Decision:        decision,
		FreeTextPresent: strings.TrimSpace(answer.FreeTextAnswer) != "",
	}); err != nil {
		return nil, err
	}

	return &SubmitAnswerResult{
		Answer:            answer,
		Package:           pkg,
		Governance:        decision,
		ProgressState:     progressState,
		AnsweredQuestions: answeredQuestions,
		TotalQuestions:    totalQuestions,
		HasMore:           answeredQuestions < totalQuestions,
	}, nil
}
