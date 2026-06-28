// Package ports declares the interfaces the sales module's services depend
// on. They are consumer-defined contracts, implemented by the adapters.
package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/example/shop/internal/modules/sales/internal/domain"
)

// OrderItemRepository is the persistence contract for order_items.
type OrderItemRepository interface {
	Create(ctx context.Context, in domain.CreateOrderItemInput) (domain.OrderItem, error)
	Get(ctx context.Context, id uuid.UUID) (domain.OrderItem, error)
	List(ctx context.Context, limit, offset int32) ([]domain.OrderItem, error)
	Count(ctx context.Context) (int64, error)
	Update(ctx context.Context, id uuid.UUID, in domain.UpdateOrderItemInput) (domain.OrderItem, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
