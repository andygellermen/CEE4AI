package http

import (
	"errors"
	"net/http"

	"github.com/andygellermen/CEE4AI/internal/reviews"
)

type reviewsHandler struct {
	service *reviews.Service
}

type governanceResponse struct {
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

type flagQuestionRequest struct {
	QuestionID   int64   `json:"question_id"`
	SessionID    *string `json:"session_id"`
	FlagSlug     string  `json:"flag_slug"`
	ReviewerRole string  `json:"reviewer_role"`
	Comment      string  `json:"comment"`
	Severity     string  `json:"severity"`
}

type questionReviewResponse struct {
	ID           int64   `json:"id"`
	QuestionID   int64   `json:"question_id"`
	SessionID    *string `json:"session_id,omitempty"`
	ReviewerRole string  `json:"reviewer_role,omitempty"`
	FlagID       int64   `json:"flag_id"`
	FlagSlug     string  `json:"flag_slug"`
	Comment      string  `json:"comment,omitempty"`
	Severity     string  `json:"severity,omitempty"`
}

type questionDecisionRequest struct {
	QuestionID int64  `json:"question_id"`
	NewStatus  string `json:"new_status"`
	Reason     string `json:"reason"`
}

type questionDecisionResponse struct {
	ID         int64  `json:"id"`
	QuestionID int64  `json:"question_id"`
	OldStatus  string `json:"old_status,omitempty"`
	NewStatus  string `json:"new_status"`
	Reason     string `json:"reason,omitempty"`
}

type localizationReviewRequest struct {
	TranslationID int64  `json:"translation_id"`
	Status        string `json:"status"`
	Comment       string `json:"comment"`
}

type localizationReviewResponse struct {
	ID            int64  `json:"id"`
	TranslationID int64  `json:"translation_id"`
	Status        string `json:"status"`
	Comment       string `json:"comment,omitempty"`
}

func newReviewsHandler(service *reviews.Service) *reviewsHandler {
	return &reviewsHandler{service: service}
}

func (h *reviewsHandler) flagQuestion(w http.ResponseWriter, r *http.Request) {
	var req flagQuestionRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.QuestionID <= 0 || req.FlagSlug == "" {
		writeError(w, http.StatusBadRequest, "question_id and flag_slug are required")
		return
	}

	review, err := h.service.FlagQuestion(r.Context(), reviews.FlagQuestionRequest{
		QuestionID:   req.QuestionID,
		SessionID:    req.SessionID,
		FlagSlug:     req.FlagSlug,
		ReviewerRole: req.ReviewerRole,
		Comment:      req.Comment,
		Severity:     req.Severity,
	})
	if err != nil {
		switch {
		case errors.Is(err, reviews.ErrReviewFlagNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusCreated, questionReviewResponse{
		ID:           review.ID,
		QuestionID:   review.QuestionID,
		SessionID:    review.SessionID,
		ReviewerRole: review.ReviewerRole,
		FlagID:       review.FlagID,
		FlagSlug:     review.FlagSlug,
		Comment:      review.Comment,
		Severity:     review.Severity,
	})
}

func (h *reviewsHandler) applyQuestionDecision(w http.ResponseWriter, r *http.Request) {
	var req questionDecisionRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.QuestionID <= 0 || req.NewStatus == "" {
		writeError(w, http.StatusBadRequest, "question_id and new_status are required")
		return
	}

	decision, err := h.service.ApplyQuestionDecision(r.Context(), reviews.ApplyQuestionDecisionRequest{
		QuestionID: req.QuestionID,
		NewStatus:  req.NewStatus,
		Reason:     req.Reason,
	})
	if err != nil {
		switch {
		case errors.Is(err, reviews.ErrQuestionNotFound), errors.Is(err, reviews.ErrReviewFlagNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, reviews.ErrInvalidReviewStatus):
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, questionDecisionResponse{
		ID:         decision.ID,
		QuestionID: decision.QuestionID,
		OldStatus:  decision.OldStatus,
		NewStatus:  decision.NewStatus,
		Reason:     decision.Reason,
	})
}

func (h *reviewsHandler) recordLocalizationReview(w http.ResponseWriter, r *http.Request) {
	var req localizationReviewRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.TranslationID <= 0 || req.Status == "" {
		writeError(w, http.StatusBadRequest, "translation_id and status are required")
		return
	}

	review, err := h.service.RecordLocalizationReview(r.Context(), reviews.RecordLocalizationReviewRequest{
		TranslationID: req.TranslationID,
		Status:        req.Status,
		Comment:       req.Comment,
	})
	if err != nil {
		switch {
		case errors.Is(err, reviews.ErrTranslationNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, reviews.ErrInvalidReviewStatus):
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusCreated, localizationReviewResponse{
		ID:            review.ID,
		TranslationID: review.TranslationID,
		Status:        review.Status,
		Comment:       review.Comment,
	})
}
