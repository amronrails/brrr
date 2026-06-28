// Package postgres implements the sales module's ports on top of the
// centralized sqlc-generated query layer (internal/db).
package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/example/shop/internal/db"
	"github.com/example/shop/internal/modules/sales/internal/domain"
	"github.com/example/shop/internal/modules/sales/internal/ports"
)

type orderItemRepository struct {
	q *db.Queries
}

// NewOrderItemRepository constructs a OrderItemRepository.
func NewOrderItemRepository(q *db.Queries) ports.OrderItemRepository {
	return &orderItemRepository{q: q}
}

func (r *orderItemRepository) Create(ctx context.Context, in domain.CreateOrderItemInput) (domain.OrderItem, error) {
	row, err := r.q.CreateOrderItem(ctx, db.CreateOrderItemParams{
		Quantity: in.Quantity,
		UnitPrice: in.UnitPrice,
		OrderID: in.OrderID,
		ProductID: in.ProductID,
	})
	if err != nil {
		return domain.OrderItem{}, mapOrderItemError(err)
	}
	return toOrderItem(row), nil
}

func (r *orderItemRepository) Get(ctx context.Context, id uuid.UUID) (domain.OrderItem, error) {
	row, err := r.q.GetOrderItem(ctx, id)
	if err != nil {
		return domain.OrderItem{}, mapOrderItemError(err)
	}
	return toOrderItem(row), nil
}

func (r *orderItemRepository) List(ctx context.Context, limit, offset int32) ([]domain.OrderItem, error) {
	rows, err := r.q.ListOrderItems(ctx, db.ListOrderItemsParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	out := make([]domain.OrderItem, 0, len(rows))
	for _, row := range rows {
		out = append(out, toOrderItem(row))
	}
	return out, nil
}

func (r *orderItemRepository) Count(ctx context.Context) (int64, error) {
	return r.q.CountOrderItems(ctx)
}

func (r *orderItemRepository) Update(ctx context.Context, id uuid.UUID, in domain.UpdateOrderItemInput) (domain.OrderItem, error) {
	row, err := r.q.UpdateOrderItem(ctx, db.UpdateOrderItemParams{
		ID: id,
		Quantity: in.Quantity,
		UnitPrice: in.UnitPrice,
		OrderID: in.OrderID,
		ProductID: in.ProductID,
	})
	if err != nil {
		return domain.OrderItem{}, mapOrderItemError(err)
	}
	return toOrderItem(row), nil
}

func (r *orderItemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteOrderItem(ctx, id)
}

func toOrderItem(row db.OrderItem) domain.OrderItem {
	return domain.OrderItem{
		ID: row.ID,
		Quantity: row.Quantity,
		UnitPrice: row.UnitPrice,
		OrderID: row.OrderID,
		ProductID: row.ProductID,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

// mapOrderItemError translates database errors into domain errors: a missing
// row becomes ErrOrderItemNotFound and a unique-constraint violation (SQLSTATE
// 23505) becomes ErrOrderItemConflict; anything else is returned unchanged.
func mapOrderItemError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrOrderItemNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return domain.ErrOrderItemConflict
	}
	return err
}
