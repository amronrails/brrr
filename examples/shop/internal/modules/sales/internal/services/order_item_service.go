// Package services holds the sales module's business logic.
package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/example/shop/internal/modules/sales/internal/domain"
	"github.com/example/shop/internal/modules/sales/internal/ports"
)

// OrderItemService implements the order_item use cases.
type OrderItemService struct {
	repo ports.OrderItemRepository
}

// NewOrderItemService constructs a OrderItemService.
func NewOrderItemService(repo ports.OrderItemRepository) *OrderItemService {
	return &OrderItemService{repo: repo}
}

func (s *OrderItemService) Create(ctx context.Context, in domain.CreateOrderItemInput) (domain.OrderItem, error) {
	return s.repo.Create(ctx, in)
}

func (s *OrderItemService) Get(ctx context.Context, id uuid.UUID) (domain.OrderItem, error) {
	return s.repo.Get(ctx, id)
}

// List returns a page of order_items along with the total count.
func (s *OrderItemService) List(ctx context.Context, limit, offset int32) ([]domain.OrderItem, int64, error) {
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

func (s *OrderItemService) Update(ctx context.Context, id uuid.UUID, in domain.UpdateOrderItemInput) (domain.OrderItem, error) {
	return s.repo.Update(ctx, id, in)
}

func (s *OrderItemService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
