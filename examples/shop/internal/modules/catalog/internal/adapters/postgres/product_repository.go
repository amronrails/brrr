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

type productRepository struct {
	q *db.Queries
}

// NewProductRepository constructs a ProductRepository.
func NewProductRepository(q *db.Queries) ports.ProductRepository {
	return &productRepository{q: q}
}

func (r *productRepository) Create(ctx context.Context, in domain.CreateProductInput) (domain.Product, error) {
	row, err := r.q.CreateProduct(ctx, db.CreateProductParams{
		Name: in.Name,
		Sku: in.Sku,
		Price: in.Price,
		Stock: in.Stock,
		Active: in.Active,
		Metadata: in.Metadata,
		CategoryID: in.CategoryID,
	})
	if err != nil {
		return domain.Product{}, mapProductError(err)
	}
	return toProduct(row), nil
}

func (r *productRepository) Get(ctx context.Context, id uuid.UUID) (domain.Product, error) {
	row, err := r.q.GetProduct(ctx, id)
	if err != nil {
		return domain.Product{}, mapProductError(err)
	}
	return toProduct(row), nil
}

func (r *productRepository) List(ctx context.Context, limit, offset int32) ([]domain.Product, error) {
	rows, err := r.q.ListProducts(ctx, db.ListProductsParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	out := make([]domain.Product, 0, len(rows))
	for _, row := range rows {
		out = append(out, toProduct(row))
	}
	return out, nil
}

func (r *productRepository) Count(ctx context.Context) (int64, error) {
	return r.q.CountProducts(ctx)
}

func (r *productRepository) Update(ctx context.Context, id uuid.UUID, in domain.UpdateProductInput) (domain.Product, error) {
	row, err := r.q.UpdateProduct(ctx, db.UpdateProductParams{
		ID: id,
		Name: in.Name,
		Sku: in.Sku,
		Price: in.Price,
		Stock: in.Stock,
		Active: in.Active,
		Metadata: in.Metadata,
		CategoryID: in.CategoryID,
	})
	if err != nil {
		return domain.Product{}, mapProductError(err)
	}
	return toProduct(row), nil
}

func (r *productRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteProduct(ctx, id)
}

func toProduct(row db.Product) domain.Product {
	return domain.Product{
		ID: row.ID,
		Name: row.Name,
		Sku: row.Sku,
		Price: row.Price,
		Stock: row.Stock,
		Active: row.Active,
		Metadata: row.Metadata,
		CategoryID: row.CategoryID,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

// mapProductError translates database errors into domain errors: a missing
// row becomes ErrProductNotFound and a unique-constraint violation (SQLSTATE
// 23505) becomes ErrProductConflict; anything else is returned unchanged.
func mapProductError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrProductNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return domain.ErrProductConflict
	}
	return err
}
