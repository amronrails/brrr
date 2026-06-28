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

type orderRepository struct {
	q *db.Queries
}

// NewOrderRepository constructs a OrderRepository.
func NewOrderRepository(q *db.Queries) ports.OrderRepository {
	return &orderRepository{q: q}
}

func (r *orderRepository) Create(ctx context.Context, in domain.CreateOrderInput) (domain.Order, error) {
	row, err := r.q.CreateOrder(ctx, db.CreateOrderParams{
		Status: in.Status,
		Total: in.Total,
		PlacedAt: in.PlacedAt,
		CustomerID: in.CustomerID,
	})
	if err != nil {
		return domain.Order{}, mapOrderError(err)
	}
	return toOrder(row), nil
}

func (r *orderRepository) Get(ctx context.Context, id uuid.UUID) (domain.Order, error) {
	row, err := r.q.GetOrder(ctx, id)
	if err != nil {
		return domain.Order{}, mapOrderError(err)
	}
	return toOrder(row), nil
}

func (r *orderRepository) List(ctx context.Context, limit, offset int32) ([]domain.Order, error) {
	rows, err := r.q.ListOrders(ctx, db.ListOrdersParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	out := make([]domain.Order, 0, len(rows))
	for _, row := range rows {
		out = append(out, toOrder(row))
	}
	return out, nil
}

func (r *orderRepository) Count(ctx context.Context) (int64, error) {
	return r.q.CountOrders(ctx)
}

func (r *orderRepository) Update(ctx context.Context, id uuid.UUID, in domain.UpdateOrderInput) (domain.Order, error) {
	row, err := r.q.UpdateOrder(ctx, db.UpdateOrderParams{
		ID: id,
		Status: in.Status,
		Total: in.Total,
		PlacedAt: in.PlacedAt,
		CustomerID: in.CustomerID,
	})
	if err != nil {
		return domain.Order{}, mapOrderError(err)
	}
	return toOrder(row), nil
}

func (r *orderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteOrder(ctx, id)
}

func toOrder(row db.Order) domain.Order {
	return domain.Order{
		ID: row.ID,
		Status: row.Status,
		Total: row.Total,
		PlacedAt: row.PlacedAt,
		CustomerID: row.CustomerID,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

// mapOrderError translates database errors into domain errors: a missing
// row becomes ErrOrderNotFound and a unique-constraint violation (SQLSTATE
// 23505) becomes ErrOrderConflict; anything else is returned unchanged.
func mapOrderError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrOrderNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return domain.ErrOrderConflict
	}
	return err
}
