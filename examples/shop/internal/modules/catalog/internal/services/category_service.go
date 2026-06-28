// Package services holds the catalog module's business logic.
package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/example/shop/internal/modules/catalog/internal/domain"
	"github.com/example/shop/internal/modules/catalog/internal/ports"
)

// CategoryService implements the category use cases.
type CategoryService struct {
	repo ports.CategoryRepository
}

// NewCategoryService constructs a CategoryService.
func NewCategoryService(repo ports.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) Create(ctx context.Context, in domain.CreateCategoryInput) (domain.Category, error) {
	return s.repo.Create(ctx, in)
}

func (s *CategoryService) Get(ctx context.Context, id uuid.UUID) (domain.Category, error) {
	return s.repo.Get(ctx, id)
}

// List returns a page of categories along with the total count.
func (s *CategoryService) List(ctx context.Context, limit, offset int32) ([]domain.Category, int64, error) {
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

func (s *CategoryService) Update(ctx context.Context, id uuid.UUID, in domain.UpdateCategoryInput) (domain.Category, error) {
	return s.repo.Update(ctx, id, in)
}

func (s *CategoryService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
