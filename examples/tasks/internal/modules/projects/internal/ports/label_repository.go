// Package ports declares the interfaces the projects module's services depend
// on. They are consumer-defined contracts, implemented by the adapters.
package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/example/tasks/internal/modules/projects/internal/domain"
)

// LabelRepository is the persistence contract for labels.
type LabelRepository interface {
	Create(ctx context.Context, in domain.CreateLabelInput) (domain.Label, error)
	Get(ctx context.Context, id uuid.UUID) (domain.Label, error)
	List(ctx context.Context, limit, offset int32) ([]domain.Label, error)
	Count(ctx context.Context) (int64, error)
	Update(ctx context.Context, id uuid.UUID, in domain.UpdateLabelInput) (domain.Label, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
