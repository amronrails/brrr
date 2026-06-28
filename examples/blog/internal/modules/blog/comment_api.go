// Comment public API DTO and accessor — the cross-module surface for
// comments. Generated once by brrr; safe to extend by hand (the aggregate
// interface lives in api.go and is regenerated).
package blog

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/example/blog/internal/modules/blog/internal/domain"
)

// Comment is the public, cross-module representation of a comment.
type Comment struct {
	ID uuid.UUID
	Body string
	Approved bool
	PostID uuid.UUID
	AuthorID uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CommentByID returns the comment with the given id.
func (m *Module) CommentByID(ctx context.Context, id uuid.UUID) (Comment, error) {
	e, err := m.commentSvc.Get(ctx, id)
	if err != nil {
		return Comment{}, err
	}
	return toAPIComment(e), nil
}

func toAPIComment(e domain.Comment) Comment {
	return Comment{
		ID: e.ID,
		Body: e.Body,
		Approved: e.Approved,
		PostID: e.PostID,
		AuthorID: e.AuthorID,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
