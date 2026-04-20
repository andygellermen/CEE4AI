package importer

type TableSpec struct {
	TableName       string
	RelativePath    string
	ConflictColumns []string
	ResetSequence   bool
}

func DefaultTableSpecs() []TableSpec {
	return []TableSpec{
		{TableName: "core.languages", RelativePath: "core/languages.csv", ConflictColumns: []string{"id"}, ResetSequence: true},
		{TableName: "core.regions", RelativePath: "core/regions.csv", ConflictColumns: []string{"id"}, ResetSequence: true},
		{TableName: "core.domains", RelativePath: "core/domains.csv", ConflictColumns: []string{"id"}, ResetSequence: true},
		{TableName: "core.roles", RelativePath: "core/roles.csv", ConflictColumns: []string{"id"}, ResetSequence: true},
		{TableName: "profiling.denktypes", RelativePath: "profiling/denktypes.csv", ConflictColumns: []string{"id"}, ResetSequence: true},
		{TableName: "profiling.skills", RelativePath: "profiling/skills.csv", ConflictColumns: []string{"id"}, ResetSequence: true},
		{TableName: "profiling.personality_traits", RelativePath: "profiling/personality_traits.csv", ConflictColumns: []string{"id"}, ResetSequence: true},
		{TableName: "profiling.interest_tags", RelativePath: "profiling/interest_tags.csv", ConflictColumns: []string{"id"}, ResetSequence: true},
		{TableName: "profiling.meaning_tags", RelativePath: "profiling/meaning_tags.csv", ConflictColumns: []string{"id"}, ResetSequence: true},
		{TableName: "profiling.worldview_frames", RelativePath: "profiling/worldview_frames.csv", ConflictColumns: []string{"id"}, ResetSequence: true},
		{TableName: "content.categories", RelativePath: "content/categories.csv", ConflictColumns: []string{"id"}, ResetSequence: true},
		{TableName: "content.subcategories", RelativePath: "content/subcategories.csv", ConflictColumns: []string{"id"}, ResetSequence: true},
		{TableName: "content.question_master", RelativePath: "content/question_master.csv", ConflictColumns: []string{"id"}, ResetSequence: true},
		{TableName: "content.question_translation", RelativePath: "content/question_translation.csv", ConflictColumns: []string{"id"}, ResetSequence: true},
		{TableName: "content.question_option_master", RelativePath: "content/question_option_master.csv", ConflictColumns: []string{"id"}, ResetSequence: true},
		{TableName: "content.question_option_translation", RelativePath: "content/question_option_translation.csv", ConflictColumns: []string{"id"}, ResetSequence: true},
		{TableName: "profiling.question_denktype_tags", RelativePath: "profiling/question_denktype_tags.csv", ConflictColumns: []string{"question_id", "denktype_id"}},
		{TableName: "profiling.question_skill_tags", RelativePath: "profiling/question_skill_tags.csv", ConflictColumns: []string{"question_id", "skill_id"}},
		{TableName: "profiling.question_trait_tags", RelativePath: "profiling/question_trait_tags.csv", ConflictColumns: []string{"question_id", "trait_id"}},
		{TableName: "profiling.question_meaning_tags", RelativePath: "profiling/question_meaning_tags.csv", ConflictColumns: []string{"question_id", "meaning_tag_id"}},
		{TableName: "profiling.question_worldview_tags", RelativePath: "profiling/question_worldview_tags.csv", ConflictColumns: []string{"question_id", "worldview_frame_id"}},
	}
}
