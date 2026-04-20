package scoring

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrInvalidAnswerInput = errors.New("invalid answer input")

type Service struct {
	pool *pgxpool.Pool
}

type EvaluateAnswerParams struct {
	QuestionID        int64
	QuestionType      string
	SelectedOptionIDs []int64
	ScaleValue        *int
	FreeTextAnswer    string
}

type Evaluation struct {
	AnswerKind        string
	SelectedOptionIDs []int64
	ScaleValue        *int
	FreeTextAnswer    string
	RawScore          *float64
	EvaluatedScore    *float64
}

type SnapshotVectors struct {
	Denktype  map[string]float64 `json:"denktype"`
	Skill     map[string]float64 `json:"skill"`
	Trait     map[string]float64 `json:"trait"`
	Meaning   map[string]float64 `json:"meaning"`
	Worldview map[string]float64 `json:"worldview"`
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

func (s *Service) EvaluateAnswer(ctx context.Context, params EvaluateAnswerParams) (*Evaluation, error) {
	evaluation := &Evaluation{
		AnswerKind:        params.QuestionType,
		SelectedOptionIDs: params.SelectedOptionIDs,
		ScaleValue:        params.ScaleValue,
		FreeTextAnswer:    strings.TrimSpace(params.FreeTextAnswer),
	}

	switch params.QuestionType {
	case "single_choice":
		if len(params.SelectedOptionIDs) != 1 {
			return nil, ErrInvalidAnswerInput
		}
		score, err := s.sumOptionWeights(ctx, params.QuestionID, params.SelectedOptionIDs)
		if err != nil {
			return nil, err
		}
		evaluation.RawScore = &score
		evaluation.EvaluatedScore = &score
	case "multiple_select":
		if len(params.SelectedOptionIDs) == 0 {
			return nil, ErrInvalidAnswerInput
		}
		score, err := s.sumOptionWeights(ctx, params.QuestionID, params.SelectedOptionIDs)
		if err != nil {
			return nil, err
		}
		evaluation.RawScore = &score
		evaluation.EvaluatedScore = &score
	case "scale":
		if params.ScaleValue == nil {
			return nil, ErrInvalidAnswerInput
		}
		score := float64(*params.ScaleValue)
		evaluation.RawScore = &score
		evaluation.EvaluatedScore = &score
	case "reflection":
		if evaluation.FreeTextAnswer == "" {
			return nil, ErrInvalidAnswerInput
		}
	default:
		return nil, ErrInvalidAnswerInput
	}

	return evaluation, nil
}

func (s *Service) BuildSessionVectors(ctx context.Context, sessionID string) (*SnapshotVectors, error) {
	denktype, err := s.queryVector(ctx, sessionID, "profiling.question_denktype_tags", "profiling.denktypes", "denktype_id")
	if err != nil {
		return nil, err
	}
	skill, err := s.queryVector(ctx, sessionID, "profiling.question_skill_tags", "profiling.skills", "skill_id")
	if err != nil {
		return nil, err
	}
	trait, err := s.queryVector(ctx, sessionID, "profiling.question_trait_tags", "profiling.personality_traits", "trait_id")
	if err != nil {
		return nil, err
	}
	meaning, err := s.queryVector(ctx, sessionID, "profiling.question_meaning_tags", "profiling.meaning_tags", "meaning_tag_id")
	if err != nil {
		return nil, err
	}
	worldview, err := s.queryVector(ctx, sessionID, "profiling.question_worldview_tags", "profiling.worldview_frames", "worldview_frame_id")
	if err != nil {
		return nil, err
	}

	return &SnapshotVectors{
		Denktype:  denktype,
		Skill:     skill,
		Trait:     trait,
		Meaning:   meaning,
		Worldview: worldview,
	}, nil
}

func TopSignal(values map[string]float64) string {
	var bestSlug string
	var bestScore float64
	for slug, score := range values {
		if bestSlug == "" || score > bestScore || (score == bestScore && slug < bestSlug) {
			bestSlug = slug
			bestScore = score
		}
	}
	return bestSlug
}

func MarshalSelectedOptionIDs(ids []int64) ([]byte, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	return json.Marshal(ids)
}

func (s *Service) sumOptionWeights(ctx context.Context, questionID int64, optionIDs []int64) (float64, error) {
	var matched int
	var score float64
	if err := s.pool.QueryRow(ctx, `
SELECT COUNT(*), COALESCE(SUM(score_weight), 0)
FROM content.question_option_master
WHERE question_id = $1
  AND id = ANY($2)
  AND is_active = TRUE
`, questionID, optionIDs).Scan(&matched, &score); err != nil {
		return 0, fmt.Errorf("sum option weights: %w", err)
	}

	if matched != len(optionIDs) {
		return 0, ErrInvalidAnswerInput
	}

	return score, nil
}

func (s *Service) queryVector(ctx context.Context, sessionID, mappingTable, referenceTable, referenceIDColumn string) (map[string]float64, error) {
	query := fmt.Sprintf(`
SELECT
    ref.slug,
    COALESCE(SUM(sig.signal * map.weight), 0)::float8 AS score
FROM runtime.answers a
JOIN %s map
  ON map.question_id = a.question_id
 AND map.is_active = TRUE
JOIN %s ref
  ON ref.id = map.%s
 AND ref.is_active = TRUE
JOIN LATERAL (
    SELECT COALESCE(
        a.evaluated_score,
        a.raw_score,
        CASE
            WHEN a.scale_value IS NOT NULL THEN a.scale_value::numeric
            WHEN NULLIF(BTRIM(COALESCE(a.free_text_answer, '')), '') IS NOT NULL THEN 1::numeric
            ELSE 1::numeric
        END
    ) AS signal
) sig ON TRUE
WHERE a.session_id = $1::uuid
GROUP BY ref.slug
ORDER BY score DESC, ref.slug
`, mappingTable, referenceTable, referenceIDColumn)

	rows, err := s.pool.Query(ctx, query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("query vector %s: %w", mappingTable, err)
	}
	defer rows.Close()

	values := make(map[string]float64)
	for rows.Next() {
		var slug string
		var score float64
		if err := rows.Scan(&slug, &score); err != nil {
			return nil, fmt.Errorf("scan vector %s: %w", mappingTable, err)
		}
		values[slug] = score
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate vector %s: %w", mappingTable, err)
	}

	return values, nil
}
