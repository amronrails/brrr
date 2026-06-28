// Package services holds the blog module's business logic.
package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/example/blog/internal/modules/blog/internal/domain"
	"github.com/example/blog/internal/modules/blog/internal/ports"
)

// PostService implements the post use cases.
type PostService struct {
	repo ports.PostRepository
}

// NewPostService constructs a PostService.
func NewPostService(repo ports.PostRepository) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) Create(ctx context.Context, in domain.CreatePostInput) (domain.Post, error) {
	return s.repo.Create(ctx, in)
}

func (s *PostService) Get(ctx context.Context, id uuid.UUID) (domain.Post, error) {
	return s.repo.Get(ctx, id)
}

// List returns a page of posts along with the total count.
func (s *PostService) List(ctx context.Context, limit, offset int32) ([]domain.Post, int64, error) {
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

func (s *PostService) Update(ctx context.Context, id uuid.UUID, in domain.UpdatePostInput) (domain.Post, error) {
	return s.repo.Update(ctx, id, in)
}

func (s *PostService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
