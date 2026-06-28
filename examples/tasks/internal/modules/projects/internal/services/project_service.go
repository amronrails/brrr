// Package services holds the projects module's business logic.
package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/example/tasks/internal/modules/projects/internal/domain"
	"github.com/example/tasks/internal/modules/projects/internal/ports"
)

// ProjectService implements the project use cases.
type ProjectService struct {
	repo ports.ProjectRepository
}

// NewProjectService constructs a ProjectService.
func NewProjectService(repo ports.ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) Create(ctx context.Context, in domain.CreateProjectInput) (domain.Project, error) {
	return s.repo.Create(ctx, in)
}

func (s *ProjectService) Get(ctx context.Context, id uuid.UUID) (domain.Project, error) {
	return s.repo.Get(ctx, id)
}

// List returns a page of projects along with the total count.
func (s *ProjectService) List(ctx context.Context, limit, offset int32) ([]domain.Project, int64, error) {
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

func (s *ProjectService) Update(ctx context.Context, id uuid.UUID, in domain.UpdateProjectInput) (domain.Project, error) {
	return s.repo.Update(ctx, id, in)
}

func (s *ProjectService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
