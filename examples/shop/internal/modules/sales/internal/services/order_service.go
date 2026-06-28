// Package services holds the sales module's business logic.
package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/example/shop/internal/modules/sales/internal/domain"
	"github.com/example/shop/internal/modules/sales/internal/ports"
)

// OrderService implements the order use cases.
type OrderService struct {
	repo ports.OrderRepository
}

// NewOrderService constructs a OrderService.
func NewOrderService(repo ports.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) Create(ctx context.Context, in domain.CreateOrderInput) (domain.Order, error) {
	return s.repo.Create(ctx, in)
}

func (s *OrderService) Get(ctx context.Context, id uuid.UUID) (domain.Order, error) {
	return s.repo.Get(ctx, id)
}

// List returns a page of orders along with the total count.
func (s *OrderService) List(ctx context.Context, limit, offset int32) ([]domain.Order, int64, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	items, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (s *OrderService) Update(ctx context.Context, id uuid.UUID, in domain.UpdateOrderInput) (domain.Order, error) {
	return s.repo.Update(ctx, id, in)
}

func (s *OrderService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
