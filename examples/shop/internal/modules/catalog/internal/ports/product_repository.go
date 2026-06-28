// Package ports declares the interfaces the catalog module's services depend
// on. They are consumer-defined contracts, implemented by the adapters.
package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/example/shop/internal/modules/catalog/internal/domain"
)

// ProductRepository is the persistence contract for products.
type ProductRepository interface {
	Create(ctx context.Context, in domain.CreateProductInput) (domain.Product, error)
	Get(ctx context.Context, id uuid.UUID) (domain.Product, error)
	List(ctx context.Context, limit, offset int32) ([]domain.Product, error)
	Count(ctx context.Context) (int64, error)
	Update(ctx context.Context, id uuid.UUID, in domain.UpdateProductInput) (domain.Product, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
