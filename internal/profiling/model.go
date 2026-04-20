package profiling

type Denktype struct {
	ID              int64
	Slug            string
	Name            string
	Description     string
	DevelopmentHint string
	IsActive        bool
}

type Skill struct {
	ID          int64
	Slug        string
	Name        string
	Description string
	IsActive    bool
}

type PersonalityTrait struct {
	ID            int64
	Slug          string
	Name          string
	Description   string
	PolarityModel string
	IsActive      bool
}

type InterestTag struct {
	ID          int64
	Slug        string
	Name        string
	Description string
	IsActive    bool
}

type MeaningTag struct {
	ID               int64
	Slug             string
	Name             string
	Description      string
	SensitivityLevel string
	IsActive         bool
}

type WorldviewFrame struct {
	ID               int64
	Slug             string
	Name             string
	Description      string
	SensitivityLevel string
	IsActive         bool
}
