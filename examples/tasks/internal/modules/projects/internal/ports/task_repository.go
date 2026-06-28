// Package ports declares the interfaces the projects module's services depend
// on. They are consumer-defined contracts, implemented by the adapters.
package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/example/tasks/internal/modules/projects/internal/domain"
)

// TaskRepository is the persistence contract for tasks.
type TaskRepository interface {
	Create(ctx context.Context, in domain.CreateTaskInput) (domain.Task, error)
	Get(ctx context.Context, id uuid.UUID) (domain.Task, error)
	List(ctx context.Context, limit, offset int32) ([]domain.Task, error)
	Count(ctx context.Context) (int64, error)
	Update(ctx context.Context, id uuid.UUID, in domain.UpdateTaskInput) (domain.Task, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
