// Post public API DTO and accessor — the cross-module surface for
// posts. Generated once by brrr; safe to extend by hand (the aggregate
// interface lives in api.go and is regenerated).
package blog

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/example/blog/internal/modules/blog/internal/domain"
)

// Post is the public, cross-module representation of a post.
type Post struct {
	ID uuid.UUID
	Title string
	Slug string
	Excerpt string
	Body string
	Published bool
	Views int32
	AuthorID uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

// PostByID returns the post with the given id.
func (m *Module) PostByID(ctx context.Context, id uuid.UUID) (Post, error) {
	e, err := m.postSvc.Get(ctx, id)
	if err != nil {
		return Post{}, err
	}
	return toAPIPost(e), nil
}

func toAPIPost(e domain.Post) Post {
	return Post{
		ID: e.ID,
		Title: e.Title,
		Slug: e.Slug,
		Excerpt: e.Excerpt,
		Body: e.Body,
		Published: e.Published,
		Views: e.Views,
		AuthorID: e.AuthorID,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
