package results

import (
	"encoding/json"
	"time"

	"github.com/andygellermen/CEE4AI/internal/governance"
	"github.com/andygellermen/CEE4AI/internal/scoring"
)

type Snapshot struct {
	ID              int64
	SessionID       string
	ResultType      string
	ProfileDepth    string
	CertaintyLevel  string
	SnapshotPayload json.RawMessage
	RulesetVersion  string
	CreatedAt       time.Time
}

type SnapshotPayload struct {
	SessionID         string                               `json:"session_id"`
	DomainID          int64                                `json:"domain_id"`
	Mode              string                               `json:"mode"`
	ProfileDepth      string                               `json:"profile_depth"`
	ProgressState     string                               `json:"progress_state"`
	AnsweredQuestions int                                  `json:"answered_questions"`
	TotalQuestions    int                                  `json:"total_questions"`
	CompletionRatio   float64                              `json:"completion_ratio"`
	ResultConfidence  float64                              `json:"result_confidence"`
	CertaintyLevel    string                               `json:"certainty_level"`
	Vectors           *scoring.SnapshotVectors             `json:"vectors"`
	Governance        *governance.SessionGovernanceSummary `json:"governance,omitempty"`
	TopSignals        map[string]string                    `json:"top_signals"`
}

type SnapshotResult struct {
	Snapshot *Snapshot
	Payload  *SnapshotPayload
}
