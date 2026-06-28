// Order public API DTO and accessor — the cross-module surface for
// orders. Generated once by brrr; safe to extend by hand (the aggregate
// interface lives in api.go and is regenerated).
package sales

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/example/shop/internal/modules/sales/internal/domain"
)

// Order is the public, cross-module representation of a order.
type Order struct {
	ID uuid.UUID
	Status string
	Total string
	PlacedAt time.Time
	CustomerID uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

// OrderByID returns the order with the given id.
func (m *Module) OrderByID(ctx context.Context, id uuid.UUID) (Order, error) {
	e, err := m.orderSvc.Get(ctx, id)
	if err != nil {
		return Order{}, err
	}
	return toAPIOrder(e), nil
}

func toAPIOrder(e domain.Order) Order {
	return Order{
		ID: e.ID,
		Status: e.Status,
		Total: e.Total,
		PlacedAt: e.PlacedAt,
		CustomerID: e.CustomerID,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
