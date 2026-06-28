// Package services holds the blog module's business logic.
package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/example/blog/internal/modules/blog/internal/domain"
	"github.com/example/blog/internal/modules/blog/internal/ports"
)

// CommentService implements the comment use cases.
type CommentService struct {
	repo ports.CommentRepository
}

// NewCommentService constructs a CommentService.
func NewCommentService(repo ports.CommentRepository) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) Create(ctx context.Context, in domain.CreateCommentInput) (domain.Comment, error) {
	return s.repo.Create(ctx, in)
}

func (s *CommentService) Get(ctx context.Context, id uuid.UUID) (domain.Comment, error) {
	return s.repo.Get(ctx, id)
}

// List returns a page of comments along with the total count.
func (s *CommentService) List(ctx context.Context, limit, offset int32) ([]domain.Comment, int64, error) {
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

func (s *CommentService) Update(ctx context.Context, id uuid.UUID, in domain.UpdateCommentInput) (domain.Comment, error) {
	return s.repo.Update(ctx, id, in)
}

func (s *CommentService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
