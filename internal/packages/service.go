package packages

import (
	"context"
)

const DefaultPackageSize = 3

type QuestionPlanner interface {
	BuildPackage(ctx context.Context, domainID int64, offset, limit int) (*QuestionPlan, error)
}

type Service struct {
	repo      *Repository
	questions QuestionPlanner
}

func NewService(repo *Repository, questions QuestionPlanner) *Service {
	return &Service{
		repo:      repo,
		questions: questions,
	}
}

func (s *Service) EnsureByQuestionPosition(ctx context.Context, sessionID string, domainID int64, position int) (*SessionPackage, error) {
	if position <= 0 {
		return nil, nil
	}

	packageIndex := ((position - 1) / DefaultPackageSize) + 1

	pkg, err := s.repo.GetBySessionAndIndex(ctx, sessionID, packageIndex)
	if err == nil {
		return pkg, nil
	}
	if err != ErrPackageNotFound {
		return nil, err
	}

	offset := (packageIndex - 1) * DefaultPackageSize
	plan, err := s.questions.BuildPackage(ctx, domainID, offset, DefaultPackageSize)
	if err != nil {
		return nil, err
	}
	if plan == nil || len(plan.QuestionIDs) == 0 {
		return nil, nil
	}

	var estimated *int
	if plan.EstimatedTimeSeconds > 0 {
		value := plan.EstimatedTimeSeconds
		estimated = &value
	}

	return s.repo.Create(ctx, CreatePackageParams{
		SessionID:            sessionID,
		PackageIndex:         packageIndex,
		PackageSize:          len(plan.QuestionIDs),
		EstimatedTimeSeconds: estimated,
	})
}
