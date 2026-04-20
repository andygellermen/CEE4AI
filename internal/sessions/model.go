package sessions

import "time"

type Session struct {
	ID               string
	DomainID         int64
	Mode             string
	SessionGoal      string
	LocaleLanguageID int64
	LocaleRegionID   *int64
	ResultConfidence *float64
	ProgressState    string
	StartedAt        time.Time
	FinishedAt       *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type CreateSessionParams struct {
	DomainID         int64
	Mode             string
	SessionGoal      string
	LocaleLanguageID int64
	LocaleRegionID   *int64
	ProgressState    string
}
