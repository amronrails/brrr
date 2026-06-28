// Package ports declares the interfaces the blog module's services depend
// on. They are consumer-defined contracts, implemented by the adapters.
package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/example/blog/internal/modules/blog/internal/domain"
)

// CommentRepository is the persistence contract for comments.
type CommentRepository interface {
	Create(ctx context.Context, in domain.CreateCommentInput) (domain.Comment, error)
	Get(ctx context.Context, id uuid.UUID) (domain.Comment, error)
	List(ctx context.Context, limit, offset int32) ([]domain.Comment, error)
	Count(ctx context.Context) (int64, error)
	Update(ctx context.Context, id uuid.UUID, in domain.UpdateCommentInput) (domain.Comment, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
