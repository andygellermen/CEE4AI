package localization

type Language struct {
	ID         int64
	Code       string
	Name       string
	ScriptCode string
	IsActive   bool
}

type Region struct {
	ID       int64
	Code     string
	Name     string
	IsActive bool
}

type TranslationMeta struct {
	LanguageID         int64
	RegionID           *int64
	LocalizationStatus string
	IsActive           bool
}
