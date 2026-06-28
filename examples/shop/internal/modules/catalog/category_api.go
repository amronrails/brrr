// Category public API DTO and accessor — the cross-module surface for
// categories. Generated once by brrr; safe to extend by hand (the aggregate
// interface lives in api.go and is regenerated).
package catalog

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/example/shop/internal/modules/catalog/internal/domain"
)

// Category is the public, cross-module representation of a category.
type Category struct {
	ID uuid.UUID
	Name string
	Slug string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CategoryByID returns the category with the given id.
func (m *Module) CategoryByID(ctx context.Context, id uuid.UUID) (Category, error) {
	e, err := m.categorySvc.Get(ctx, id)
	if err != nil {
		return Category{}, err
	}
	return toAPICategory(e), nil
}

func toAPICategory(e domain.Category) Category {
	return Category{
		ID: e.ID,
		Name: e.Name,
		Slug: e.Slug,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
