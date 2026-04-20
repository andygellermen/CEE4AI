package http

import (
	"errors"
	"net/http"

	"github.com/andygellermen/CEE4AI/internal/governance"
	"github.com/andygellermen/CEE4AI/internal/questions"
	"github.com/andygellermen/CEE4AI/internal/sessions"
)

type questionsHandler struct {
	service *questions.Service
}

type nextQuestionResponse struct {
	Session            sessionResponse     `json:"session"`
	Package            *packageResponse    `json:"package,omitempty"`
	Question           *questionResponse   `json:"question,omitempty"`
	Governance         *governanceResponse `json:"governance,omitempty"`
	RemainingQuestions int                 `json:"remaining_questions"`
	HasMore            bool                `json:"has_more"`
}

type questionResponse struct {
	ID                              int64            `json:"id"`
	ExternalID                      string           `json:"external_id"`
	QuestionFamily                  string           `json:"question_family"`
	QuestionType                    string           `json:"question_type"`
	IntendedUse                     string           `json:"intended_use,omitempty"`
	ConfidenceTier                  string           `json:"confidence_tier,omitempty"`
	EstimatedTimeSeconds            *int             `json:"estimated_time_seconds,omitempty"`
	CognitiveLoadLevel              string           `json:"cognitive_load_level,omitempty"`
	MeaningDepth                    string           `json:"meaning_depth,omitempty"`
	WorldviewSensitivity            string           `json:"worldview_sensitivity,omitempty"`
	SymbolicInterpretationRelevance string           `json:"symbolic_interpretation_relevance,omitempty"`
	ExistentialLoadLevel            string           `json:"existential_load_level,omitempty"`
	IsSensitive                     bool             `json:"is_sensitive"`
	AgeGate                         int              `json:"age_gate"`
	Title                           string           `json:"title,omitempty"`
	QuestionText                    string           `json:"question_text"`
	ExplanationText                 string           `json:"explanation_text,omitempty"`
	RequiresHumanReview             bool             `json:"requires_human_review"`
	WorldviewSensitive              bool             `json:"worldview_sensitive"`
	Options                         []optionResponse `json:"options,omitempty"`
}

type optionResponse struct {
	ID         int64  `json:"id"`
	OptionKey  string `json:"option_key"`
	OptionText string `json:"option_text"`
}

func newQuestionsHandler(service *questions.Service) *questionsHandler {
	return &questionsHandler{service: service}
}

func (h *questionsHandler) next(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		writeError(w, http.StatusBadRequest, "session_id is required")
		return
	}

	result, err := h.service.Next(r.Context(), sessionID)
	if err != nil {
		if errors.Is(err, sessions.ErrSessionNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := nextQuestionResponse{
		Session:            makeSessionResponse(result.Session),
		Package:            makePackageResponse(result.Package),
		Governance:         makeGovernanceResponse(result.Governance),
		RemainingQuestions: result.RemainingQuestions,
		HasMore:            result.HasMore,
	}
	if result.Question != nil {
		response.Question = makeQuestionResponse(result.Question)
	}

	writeJSON(w, http.StatusOK, response)
}

func makeQuestionResponse(question *questions.DeliveryQuestion) *questionResponse {
	options := make([]optionResponse, 0, len(question.Options))
	for _, option := range question.Options {
		options = append(options, optionResponse{
			ID:         option.ID,
			OptionKey:  option.OptionKey,
			OptionText: option.OptionText,
		})
	}

	return &questionResponse{
		ID:                              question.ID,
		ExternalID:                      question.ExternalID,
		QuestionFamily:                  question.QuestionFamily,
		QuestionType:                    question.QuestionType,
		IntendedUse:                     question.IntendedUse,
		ConfidenceTier:                  question.ConfidenceTier,
		EstimatedTimeSeconds:            question.EstimatedTimeSeconds,
		CognitiveLoadLevel:              question.CognitiveLoadLevel,
		MeaningDepth:                    question.MeaningDepth,
		WorldviewSensitivity:            question.WorldviewSensitivity,
		SymbolicInterpretationRelevance: question.SymbolicInterpretationRelevance,
		ExistentialLoadLevel:            question.ExistentialLoadLevel,
		IsSensitive:                     question.IsSensitive,
		AgeGate:                         question.AgeGate,
		Title:                           question.Title,
		QuestionText:                    question.QuestionText,
		ExplanationText:                 question.ExplanationText,
		RequiresHumanReview:             question.RequiresHumanReview,
		WorldviewSensitive:              question.WorldviewSensitive,
		Options:                         options,
	}
}

func makeGovernanceResponse(decision *governance.RuntimeDecision) *governanceResponse {
	if decision == nil {
		return nil
	}

	return &governanceResponse{
		RulesetSlug:         decision.RulesetSlug,
		RulesetVersion:      decision.RulesetVersion,
		DeliveryMode:        decision.DeliveryMode,
		Severity:            decision.Severity,
		ReviewRequired:      decision.ReviewRequired,
		ApprovedForDelivery: decision.ApprovedForDelivery,
		SensitivityFlags:    decision.SensitivityFlags,
		AppliedPolicies:     decision.AppliedPolicies,
		Notes:               decision.Notes,
	}
}
