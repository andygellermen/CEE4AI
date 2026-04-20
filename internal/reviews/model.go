package reviews

import "time"

type QuestionReview struct {
	ID           int64
	QuestionID   int64
	SessionID    *string
	ReviewerRole string
	FlagID       int64
	FlagSlug     string
	Comment      string
	Severity     string
	CreatedAt    time.Time
}

type ReviewDecision struct {
	ID         int64
	QuestionID int64
	OldStatus  string
	NewStatus  string
	Reason     string
	CreatedAt  time.Time
}

type LocalizationReview struct {
	ID            int64
	TranslationID int64
	Status        string
	Comment       string
	CreatedAt     time.Time
}

type TranslationTarget struct {
	ID                  int64
	QuestionID          int64
	LocalizationStatus  string
	RequiresHumanReview bool
	WorldviewSensitive  bool
}
