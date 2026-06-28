// Package services holds the projects module's business logic.
package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/example/tasks/internal/modules/projects/internal/domain"
	"github.com/example/tasks/internal/modules/projects/internal/ports"
)

// LabelService implements the label use cases.
type LabelService struct {
	repo ports.LabelRepository
}

// NewLabelService constructs a LabelService.
func NewLabelService(repo ports.LabelRepository) *LabelService {
	return &LabelService{repo: repo}
}

func (s *LabelService) Create(ctx context.Context, in domain.CreateLabelInput) (domain.Label, error) {
	return s.repo.Create(ctx, in)
}

func (s *LabelService) Get(ctx context.Context, id uuid.UUID) (domain.Label, error) {
	return s.repo.Get(ctx, id)
}

// List returns a page of labels along with the total count.
func (s *LabelService) List(ctx context.Context, limit, offset int32) ([]domain.Label, int64, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	items, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (s *LabelService) Update(ctx context.Context, id uuid.UUID, in domain.UpdateLabelInput) (domain.Label, error) {
	return s.repo.Update(ctx, id, in)
}

func (s *LabelService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
