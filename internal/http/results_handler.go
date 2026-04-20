package http

import (
	"errors"
	"net/http"

	"github.com/andygellermen/CEE4AI/internal/results"
	"github.com/andygellermen/CEE4AI/internal/sessions"
)

type resultsHandler struct {
	service *results.Service
}

type snapshotResponse struct {
	ID             int64                    `json:"id"`
	ResultType     string                   `json:"result_type"`
	ProfileDepth   string                   `json:"profile_depth"`
	CertaintyLevel string                   `json:"certainty_level"`
	RulesetVersion string                   `json:"ruleset_version"`
	Payload        *results.SnapshotPayload `json:"payload"`
}

func newResultsHandler(service *results.Service) *resultsHandler {
	return &resultsHandler{service: service}
}

func (h *resultsHandler) get(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		writeError(w, http.StatusBadRequest, "session_id is required")
		return
	}

	result, err := h.service.BuildSnapshot(r.Context(), sessionID)
	if err != nil {
		if errors.Is(err, sessions.ErrSessionNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, snapshotResponse{
		ID:             result.Snapshot.ID,
		ResultType:     result.Snapshot.ResultType,
		ProfileDepth:   result.Snapshot.ProfileDepth,
		CertaintyLevel: result.Snapshot.CertaintyLevel,
		RulesetVersion: result.Snapshot.RulesetVersion,
		Payload:        result.Payload,
	})
}
