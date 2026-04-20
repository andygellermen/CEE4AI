package questions

type Category struct {
	ID                   int64
	DomainID             int64
	Slug                 string
	Name                 string
	Description          string
	DefaultSensitivity   string
	CulturalScope        string
	SpiritualRelevance   string
	WorldviewSensitivity string
	MeaningPathRelevance string
	IsSensitive          bool
	IsActive             bool
}

type Subcategory struct {
	ID          int64
	CategoryID  int64
	Slug        string
	Name        string
	Description string
	IsSensitive bool
	IsActive    bool
}

type Question struct {
	ID                              int64
	ExternalID                      string
	DomainID                        int64
	QuestionFamily                  string
	CategoryID                      int64
	SubcategoryID                   *int64
	ParentQuestionID                *int64
	QuestionType                    string
	ScoringMode                     string
	IntendedUse                     string
	ConfidenceTier                  string
	EstimatedTimeSeconds            *int
	CognitiveLoadLevel              string
	CulturalScope                   string
	RegionScope                     string
	MeaningDepth                    string
	WorldviewSensitivity            string
	SymbolicInterpretationRelevance string
	ExistentialLoadLevel            string
	IsSensitive                     bool
	AgeGate                         int
	ReviewStatus                    string
	IsActive                        bool
	Version                         int
}

type Translation struct {
	ID                  int64
	QuestionID          int64
	LanguageID          int64
	RegionID            *int64
	LocalizationStatus  string
	RequiresHumanReview bool
	WorldviewSensitive  bool
	Title               string
	QuestionText        string
	ExplanationText     string
	ReviewerNotes       string
	IsActive            bool
	Version             int
}

type Option struct {
	ID           int64
	QuestionID   int64
	OptionKey    string
	ScoreWeight  float64
	IsCorrect    bool
	DisplayOrder int
	IsActive     bool
	Version      int
}

type OptionTranslation struct {
	ID                 int64
	OptionID           int64
	LanguageID         int64
	RegionID           *int64
	OptionText         string
	LocalizationStatus string
	IsActive           bool
}
