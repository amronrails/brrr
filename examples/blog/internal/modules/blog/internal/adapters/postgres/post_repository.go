// Package postgres implements the blog module's ports on top of the
// centralized sqlc-generated query layer (internal/db).
package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/example/blog/internal/db"
	"github.com/example/blog/internal/modules/blog/internal/domain"
	"github.com/example/blog/internal/modules/blog/internal/ports"
)

type postRepository struct {
	q *db.Queries
}

// NewPostRepository constructs a PostRepository.
func NewPostRepository(q *db.Queries) ports.PostRepository {
	return &postRepository{q: q}
}

func (r *postRepository) Create(ctx context.Context, in domain.CreatePostInput) (domain.Post, error) {
	row, err := r.q.CreatePost(ctx, db.CreatePostParams{
		Title: in.Title,
		Slug: in.Slug,
		Excerpt: in.Excerpt,
		Body: in.Body,
		Published: in.Published,
		Views: in.Views,
		AuthorID: in.AuthorID,
	})
	if err != nil {
		return domain.Post{}, mapPostError(err)
	}
	return toPost(row), nil
}

func (r *postRepository) Get(ctx context.Context, id uuid.UUID) (domain.Post, error) {
	row, err := r.q.GetPost(ctx, id)
	if err != nil {
		return domain.Post{}, mapPostError(err)
	}
	return toPost(row), nil
}

func (r *postRepository) List(ctx context.Context, limit, offset int32) ([]domain.Post, error) {
	rows, err := r.q.ListPosts(ctx, db.ListPostsParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	out := make([]domain.Post, 0, len(rows))
	for _, row := range rows {
		out = append(out, toPost(row))
	}
	return out, nil
}

func (r *postRepository) Count(ctx context.Context) (int64, error) {
	return r.q.CountPosts(ctx)
}

func (r *postRepository) Update(ctx context.Context, id uuid.UUID, in domain.UpdatePostInput) (domain.Post, error) {
	row, err := r.q.UpdatePost(ctx, db.UpdatePostParams{
		ID: id,
		Title: in.Title,
		Slug: in.Slug,
		Excerpt: in.Excerpt,
		Body: in.Body,
		Published: in.Published,
		Views: in.Views,
		AuthorID: in.AuthorID,
	})
	if err != nil {
		return domain.Post{}, mapPostError(err)
	}
	return toPost(row), nil
}

func (r *postRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.q.DeletePost(ctx, id)
}

func toPost(row db.Post) domain.Post {
	return domain.Post{
		ID: row.ID,
		Title: row.Title,
		Slug: row.Slug,
		Excerpt: row.Excerpt,
		Body: row.Body,
		Published: row.Published,
		Views: row.Views,
		AuthorID: row.AuthorID,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

// mapPostError translates database errors into domain errors: a missing
// row becomes ErrPostNotFound and a unique-constraint violation (SQLSTATE
// 23505) becomes ErrPostConflict; anything else is returned unchanged.
func mapPostError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrPostNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return domain.ErrPostConflict
	}
	return err
}
