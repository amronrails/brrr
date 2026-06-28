// OrderItem public API DTO and accessor — the cross-module surface for
// order_items. Generated once by brrr; safe to extend by hand (the aggregate
// interface lives in api.go and is regenerated).
package sales

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/example/shop/internal/modules/sales/internal/domain"
)

// OrderItem is the public, cross-module representation of a order_item.
type OrderItem struct {
	ID uuid.UUID
	Quantity int32
	UnitPrice string
	OrderID uuid.UUID
	ProductID uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

// OrderItemByID returns the order_item with the given id.
func (m *Module) OrderItemByID(ctx context.Context, id uuid.UUID) (OrderItem, error) {
	e, err := m.orderItemSvc.Get(ctx, id)
	if err != nil {
		return OrderItem{}, err
	}
	return toAPIOrderItem(e), nil
}

func toAPIOrderItem(e domain.OrderItem) OrderItem {
	return OrderItem{
		ID: e.ID,
		Quantity: e.Quantity,
		UnitPrice: e.UnitPrice,
		OrderID: e.OrderID,
		ProductID: e.ProductID,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
