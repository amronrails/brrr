// Package ports declares the interfaces the sales module's services depend
// on. They are consumer-defined contracts, implemented by the adapters.
package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/example/shop/internal/modules/sales/internal/domain"
)

// OrderRepository is the persistence contract for orders.
type OrderRepository interface {
	Create(ctx context.Context, in domain.CreateOrderInput) (domain.Order, error)
	Get(ctx context.Context, id uuid.UUID) (domain.Order, error)
	List(ctx context.Context, limit, offset int32) ([]domain.Order, error)
	Count(ctx context.Context) (int64, error)
	Update(ctx context.Context, id uuid.UUID, in domain.UpdateOrderInput) (domain.Order, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
