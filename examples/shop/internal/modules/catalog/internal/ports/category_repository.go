// Package ports declares the interfaces the catalog module's services depend
// on. They are consumer-defined contracts, implemented by the adapters.
package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/example/shop/internal/modules/catalog/internal/domain"
)

// CategoryRepository is the persistence contract for categories.
type CategoryRepository interface {
	Create(ctx context.Context, in domain.CreateCategoryInput) (domain.Category, error)
	Get(ctx context.Context, id uuid.UUID) (domain.Category, error)
	List(ctx context.Context, limit, offset int32) ([]domain.Category, error)
	Count(ctx context.Context) (int64, error)
	Update(ctx context.Context, id uuid.UUID, in domain.UpdateCategoryInput) (domain.Category, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
