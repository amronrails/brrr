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

type commentRepository struct {
	q *db.Queries
}

// NewCommentRepository constructs a CommentRepository.
func NewCommentRepository(q *db.Queries) ports.CommentRepository {
	return &commentRepository{q: q}
}

func (r *commentRepository) Create(ctx context.Context, in domain.CreateCommentInput) (domain.Comment, error) {
	row, err := r.q.CreateComment(ctx, db.CreateCommentParams{
		Body: in.Body,
		Approved: in.Approved,
		PostID: in.PostID,
		AuthorID: in.AuthorID,
	})
	if err != nil {
		return domain.Comment{}, mapCommentError(err)
	}
	return toComment(row), nil
}

func (r *commentRepository) Get(ctx context.Context, id uuid.UUID) (domain.Comment, error) {
	row, err := r.q.GetComment(ctx, id)
	if err != nil {
		return domain.Comment{}, mapCommentError(err)
	}
	return toComment(row), nil
}

func (r *commentRepository) List(ctx context.Context, limit, offset int32) ([]domain.Comment, error) {
	rows, err := r.q.ListComments(ctx, db.ListCommentsParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	out := make([]domain.Comment, 0, len(rows))
	for _, row := range rows {
		out = append(out, toComment(row))
	}
	return out, nil
}

func (r *commentRepository) Count(ctx context.Context) (int64, error) {
	return r.q.CountComments(ctx)
}

func (r *commentRepository) Update(ctx context.Context, id uuid.UUID, in domain.UpdateCommentInput) (domain.Comment, error) {
	row, err := r.q.UpdateComment(ctx, db.UpdateCommentParams{
		ID: id,
		Body: in.Body,
		Approved: in.Approved,
		PostID: in.PostID,
		AuthorID: in.AuthorID,
	})
	if err != nil {
		return domain.Comment{}, mapCommentError(err)
	}
	return toComment(row), nil
}

func (r *commentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteComment(ctx, id)
}

func toComment(row db.Comment) domain.Comment {
	return domain.Comment{
		ID: row.ID,
		Body: row.Body,
		Approved: row.Approved,
		PostID: row.PostID,
		AuthorID: row.AuthorID,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

// mapCommentError translates database errors into domain errors: a missing
// row becomes ErrCommentNotFound and a unique-constraint violation (SQLSTATE
// 23505) becomes ErrCommentConflict; anything else is returned unchanged.
func mapCommentError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrCommentNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return domain.ErrCommentConflict
	}
	return err
}
