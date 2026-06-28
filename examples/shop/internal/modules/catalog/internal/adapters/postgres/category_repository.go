// Package postgres implements the catalog module's ports on top of the
// centralized sqlc-generated query layer (internal/db).
package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/example/shop/internal/db"
	"github.com/example/shop/internal/modules/catalog/internal/domain"
	"github.com/example/shop/internal/modules/catalog/internal/ports"
)

type categoryRepository struct {
	q *db.Queries
}

// NewCategoryRepository constructs a CategoryRepository.
func NewCategoryRepository(q *db.Queries) ports.CategoryRepository {
	return &categoryRepository{q: q}
}

func (r *categoryRepository) Create(ctx context.Context, in domain.CreateCategoryInput) (domain.Category, error) {
	row, err := r.q.CreateCategory(ctx, db.CreateCategoryParams{
		Name: in.Name,
		Slug: in.Slug,
	})
	if err != nil {
		return domain.Category{}, mapCategoryError(err)
	}
	return toCategory(row), nil
}

func (r *categoryRepository) Get(ctx context.Context, id uuid.UUID) (domain.Category, error) {
	row, err := r.q.GetCategory(ctx, id)
	if err != nil {
		return domain.Category{}, mapCategoryError(err)
	}
	return toCategory(row), nil
}

func (r *categoryRepository) List(ctx context.Context, limit, offset int32) ([]domain.Category, error) {
	rows, err := r.q.ListCategories(ctx, db.ListCategoriesParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	out := make([]domain.Category, 0, len(rows))
	for _, row := range rows {
		out = append(out, toCategory(row))
	}
	return out, nil
}

func (r *categoryRepository) Count(ctx context.Context) (int64, error) {
	return r.q.CountCategories(ctx)
}

func (r *categoryRepository) Update(ctx context.Context, id uuid.UUID, in domain.UpdateCategoryInput) (domain.Category, error) {
	row, err := r.q.UpdateCategory(ctx, db.UpdateCategoryParams{
		ID: id,
		Name: in.Name,
		Slug: in.Slug,
	})
	if err != nil {
		return domain.Category{}, mapCategoryError(err)
	}
	return toCategory(row), nil
}

func (r *categoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteCategory(ctx, id)
}

func toCategory(row db.Category) domain.Category {
	return domain.Category{
		ID: row.ID,
		Name: row.Name,
		Slug: row.Slug,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

// mapCategoryError translates database errors into domain errors: a missing
// row becomes ErrCategoryNotFound and a unique-constraint violation (SQLSTATE
// 23505) becomes ErrCategoryConflict; anything else is returned unchanged.
func mapCategoryError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrCategoryNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return domain.ErrCategoryConflict
	}
	return err
}
