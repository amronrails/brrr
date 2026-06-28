// Package services holds the catalog module's business logic.
package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/example/shop/internal/modules/catalog/internal/domain"
	"github.com/example/shop/internal/modules/catalog/internal/ports"
)

// ProductService implements the product use cases.
type ProductService struct {
	repo ports.ProductRepository
}

// NewProductService constructs a ProductService.
func NewProductService(repo ports.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) Create(ctx context.Context, in domain.CreateProductInput) (domain.Product, error) {
	return s.repo.Create(ctx, in)
}

func (s *ProductService) Get(ctx context.Context, id uuid.UUID) (domain.Product, error) {
	return s.repo.Get(ctx, id)
}

// List returns a page of products along with the total count.
func (s *ProductService) List(ctx context.Context, limit, offset int32) ([]domain.Product, int64, error) {
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

func (s *ProductService) Update(ctx context.Context, id uuid.UUID, in domain.UpdateProductInput) (domain.Product, error) {
	return s.repo.Update(ctx, id, in)
}

func (s *ProductService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
