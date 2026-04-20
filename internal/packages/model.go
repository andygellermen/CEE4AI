package packages

import "time"

type SessionPackage struct {
	ID                      int64
	SessionID               string
	PackageIndex            int
	PackageSize             int
	EstimatedTimeSeconds    *int
	ActualTimeSeconds       *int
	CompletionQuality       *float64
	ContinuationWindowUntil *time.Time
	RecommendedNextMode     *string
	CreatedAt               time.Time
}

type CreatePackageParams struct {
	SessionID               string
	PackageIndex            int
	PackageSize             int
	EstimatedTimeSeconds    *int
	RecommendedNextMode     *string
	ContinuationWindowUntil *time.Time
}

type QuestionPlan struct {
	QuestionIDs          []int64
	EstimatedTimeSeconds int
}
