package answers

import "time"

type Answer struct {
	ID                int64
	SessionID         string
	PackageID         *int64
	QuestionID        int64
	AnswerKind        string
	SelectedOptionIDs []int64
	ScaleValue        *int
	FreeTextAnswer    string
	RawScore          *float64
	EvaluatedScore    *float64
	CertaintyLevel    string
	AnsweredAt        time.Time
	CreatedAt         time.Time
}

type CreateAnswerParams struct {
	SessionID         string
	PackageID         *int64
	QuestionID        int64
	AnswerKind        string
	SelectedOptionIDs []byte
	ScaleValue        *int
	FreeTextAnswer    string
	RawScore          *float64
	EvaluatedScore    *float64
	CertaintyLevel    string
}
