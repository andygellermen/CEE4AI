package http

import (
	"errors"
	"net/http"

	"github.com/andygellermen/CEE4AI/internal/packages"
	"github.com/andygellermen/CEE4AI/internal/sessions"
)

type sessionHandler struct {
	service *sessions.Service
}

type createSessionRequest struct {
	DomainID         int64  `json:"domain_id"`
	Mode             string `json:"mode"`
	SessionGoal      string `json:"session_goal"`
	LocaleLanguageID *int64 `json:"locale_language_id"`
	LocaleRegionID   *int64 `json:"locale_region_id"`
}

type sessionResponse struct {
	ID               string   `json:"id"`
	DomainID         int64    `json:"domain_id"`
	Mode             string   `json:"mode"`
	SessionGoal      string   `json:"session_goal,omitempty"`
	LocaleLanguageID int64    `json:"locale_language_id"`
	LocaleRegionID   *int64   `json:"locale_region_id,omitempty"`
	ResultConfidence *float64 `json:"result_confidence,omitempty"`
	ProgressState    string   `json:"progress_state"`
}

type packageResponse struct {
	ID                   int64   `json:"id"`
	PackageIndex         int     `json:"package_index"`
	PackageSize          int     `json:"package_size"`
	EstimatedTimeSeconds *int    `json:"estimated_time_seconds,omitempty"`
	RecommendedNextMode  *string `json:"recommended_next_mode,omitempty"`
}

type createSessionResponse struct {
	Session      sessionResponse  `json:"session"`
	FirstPackage *packageResponse `json:"first_package,omitempty"`
}

func newSessionHandler(service *sessions.Service) *sessionHandler {
	return &sessionHandler{service: service}
}

func (h *sessionHandler) create(w http.ResponseWriter, r *http.Request) {
	var req createSessionRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.DomainID <= 0 {
		writeError(w, http.StatusBadRequest, "domain_id is required")
		return
	}

	result, err := h.service.Start(r.Context(), sessions.StartSessionRequest{
		DomainID:         req.DomainID,
		Mode:             req.Mode,
		SessionGoal:      req.SessionGoal,
		LocaleLanguageID: req.LocaleLanguageID,
		LocaleRegionID:   req.LocaleRegionID,
	})
	if err != nil {
		if errors.Is(err, sessions.ErrInvalidMode) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, createSessionResponse{
		Session:      makeSessionResponse(result.Session),
		FirstPackage: makePackageResponse(result.FirstPackage),
	})
}

func makeSessionResponse(session *sessions.Session) sessionResponse {
	return sessionResponse{
		ID:               session.ID,
		DomainID:         session.DomainID,
		Mode:             session.Mode,
		SessionGoal:      session.SessionGoal,
		LocaleLanguageID: session.LocaleLanguageID,
		LocaleRegionID:   session.LocaleRegionID,
		ResultConfidence: session.ResultConfidence,
		ProgressState:    session.ProgressState,
	}
}

func makePackageResponse(pkg *packages.SessionPackage) *packageResponse {
	if pkg == nil {
		return nil
	}

	return &packageResponse{
		ID:                   pkg.ID,
		PackageIndex:         pkg.PackageIndex,
		PackageSize:          pkg.PackageSize,
		EstimatedTimeSeconds: pkg.EstimatedTimeSeconds,
		RecommendedNextMode:  pkg.RecommendedNextMode,
	}
}
