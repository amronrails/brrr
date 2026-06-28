// Package ports declares the interfaces the projects module's services depend
// on. They are consumer-defined contracts, implemented by the adapters.
package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/example/tasks/internal/modules/projects/internal/domain"
)

// ProjectRepository is the persistence contract for projects.
type ProjectRepository interface {
	Create(ctx context.Context, in domain.CreateProjectInput) (domain.Project, error)
	Get(ctx context.Context, id uuid.UUID) (domain.Project, error)
	List(ctx context.Context, limit, offset int32) ([]domain.Project, error)
	Count(ctx context.Context) (int64, error)
	Update(ctx context.Context, id uuid.UUID, in domain.UpdateProjectInput) (domain.Project, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
