// Package ports declares the interfaces the blog module's services depend
// on. They are consumer-defined contracts, implemented by the adapters.
package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/example/blog/internal/modules/blog/internal/domain"
)

// PostRepository is the persistence contract for posts.
type PostRepository interface {
	Create(ctx context.Context, in domain.CreatePostInput) (domain.Post, error)
	Get(ctx context.Context, id uuid.UUID) (domain.Post, error)
	List(ctx context.Context, limit, offset int32) ([]domain.Post, error)
	Count(ctx context.Context) (int64, error)
	Update(ctx context.Context, id uuid.UUID, in domain.UpdatePostInput) (domain.Post, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
