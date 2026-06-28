// Package services holds the projects module's business logic.
package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/example/tasks/internal/modules/projects/internal/domain"
	"github.com/example/tasks/internal/modules/projects/internal/ports"
)

// TaskService implements the task use cases.
type TaskService struct {
	repo ports.TaskRepository
}

// NewTaskService constructs a TaskService.
func NewTaskService(repo ports.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) Create(ctx context.Context, in domain.CreateTaskInput) (domain.Task, error) {
	return s.repo.Create(ctx, in)
}

func (s *TaskService) Get(ctx context.Context, id uuid.UUID) (domain.Task, error) {
	return s.repo.Get(ctx, id)
}

// List returns a page of tasks along with the total count.
func (s *TaskService) List(ctx context.Context, limit, offset int32) ([]domain.Task, int64, error) {
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

func (s *TaskService) Update(ctx context.Context, id uuid.UUID, in domain.UpdateTaskInput) (domain.Task, error) {
	return s.repo.Update(ctx, id, in)
}

func (s *TaskService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
