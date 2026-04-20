package http

import (
	"errors"
	"net/http"

	"github.com/andygellermen/CEE4AI/internal/answers"
	"github.com/andygellermen/CEE4AI/internal/scoring"
	"github.com/andygellermen/CEE4AI/internal/sessions"
)

type answersHandler struct {
	service *answers.Service
}

type submitAnswerRequest struct {
	SessionID         string  `json:"session_id"`
	QuestionID        int64   `json:"question_id"`
	SelectedOptionIDs []int64 `json:"selected_option_ids"`
	ScaleValue        *int    `json:"scale_value"`
	FreeTextAnswer    string  `json:"free_text_answer"`
	CertaintyLevel    string  `json:"certainty_level"`
}

type answerResponse struct {
	ID                int64    `json:"id"`
	SessionID         string   `json:"session_id"`
	PackageID         *int64   `json:"package_id,omitempty"`
	QuestionID        int64    `json:"question_id"`
	AnswerKind        string   `json:"answer_kind"`
	SelectedOptionIDs []int64  `json:"selected_option_ids,omitempty"`
	ScaleValue        *int     `json:"scale_value,omitempty"`
	FreeTextAnswer    string   `json:"free_text_answer,omitempty"`
	RawScore          *float64 `json:"raw_score,omitempty"`
	EvaluatedScore    *float64 `json:"evaluated_score,omitempty"`
	CertaintyLevel    string   `json:"certainty_level,omitempty"`
}

type submitAnswerResponse struct {
	Answer            answerResponse      `json:"answer"`
	Package           *packageResponse    `json:"package,omitempty"`
	Governance        *governanceResponse `json:"governance,omitempty"`
	ProgressState     string              `json:"progress_state"`
	AnsweredQuestions int                 `json:"answered_questions"`
	TotalQuestions    int                 `json:"total_questions"`
	HasMore           bool                `json:"has_more"`
}

func newAnswersHandler(service *answers.Service) *answersHandler {
	return &answersHandler{service: service}
}

func (h *answersHandler) create(w http.ResponseWriter, r *http.Request) {
	var req submitAnswerRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.SessionID == "" || req.QuestionID <= 0 {
		writeError(w, http.StatusBadRequest, "session_id and question_id are required")
		return
	}

	result, err := h.service.Submit(r.Context(), answers.SubmitAnswerRequest{
		SessionID:         req.SessionID,
		QuestionID:        req.QuestionID,
		SelectedOptionIDs: req.SelectedOptionIDs,
		ScaleValue:        req.ScaleValue,
		FreeTextAnswer:    req.FreeTextAnswer,
		CertaintyLevel:    req.CertaintyLevel,
	})
	if err != nil {
		switch {
		case errors.Is(err, sessions.ErrSessionNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, answers.ErrQuestionAlreadyAnswered):
			writeError(w, http.StatusConflict, err.Error())
		case errors.Is(err, answers.ErrQuestionSessionMismatch), errors.Is(err, scoring.ErrInvalidAnswerInput):
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusCreated, submitAnswerResponse{
		Answer:            makeAnswerResponse(result.Answer),
		Package:           makePackageResponse(result.Package),
		Governance:        makeGovernanceResponse(result.Governance),
		ProgressState:     result.ProgressState,
		AnsweredQuestions: result.AnsweredQuestions,
		TotalQuestions:    result.TotalQuestions,
		HasMore:           result.HasMore,
	})
}

func makeAnswerResponse(answer *answers.Answer) answerResponse {
	return answerResponse{
		ID:                answer.ID,
		SessionID:         answer.SessionID,
		PackageID:         answer.PackageID,
		QuestionID:        answer.QuestionID,
		AnswerKind:        answer.AnswerKind,
		SelectedOptionIDs: answer.SelectedOptionIDs,
		ScaleValue:        answer.ScaleValue,
		FreeTextAnswer:    answer.FreeTextAnswer,
		RawScore:          answer.RawScore,
		EvaluatedScore:    answer.EvaluatedScore,
		CertaintyLevel:    answer.CertaintyLevel,
	}
}
