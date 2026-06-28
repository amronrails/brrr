// Product public API DTO and accessor — the cross-module surface for
// products. Generated once by brrr; safe to extend by hand (the aggregate
// interface lives in api.go and is regenerated).
package catalog

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/example/shop/internal/modules/catalog/internal/domain"
)

// Product is the public, cross-module representation of a product.
type Product struct {
	ID uuid.UUID
	Name string
	Sku string
	Price string
	Stock int32
	Active bool
	Metadata json.RawMessage
	CategoryID uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ProductByID returns the product with the given id.
func (m *Module) ProductByID(ctx context.Context, id uuid.UUID) (Product, error) {
	e, err := m.productSvc.Get(ctx, id)
	if err != nil {
		return Product{}, err
	}
	return toAPIProduct(e), nil
}

func toAPIProduct(e domain.Product) Product {
	return Product{
		ID: e.ID,
		Name: e.Name,
		Sku: e.Sku,
		Price: e.Price,
		Stock: e.Stock,
		Active: e.Active,
		Metadata: e.Metadata,
		CategoryID: e.CategoryID,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
