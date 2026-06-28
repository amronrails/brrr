// Package domain holds the catalog module's entities and value objects.
package domain

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Product is the product entity.
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

// CreateProductInput carries the data needed to create a product.
type CreateProductInput struct {
	Name string
	Sku string
	Price string
	Stock int32
	Active bool
	Metadata json.RawMessage
	CategoryID uuid.UUID
}

// UpdateProductInput carries the data needed to update a product.
type UpdateProductInput struct {
	Name string
	Sku string
	Price string
	Stock int32
	Active bool
	Metadata json.RawMessage
	CategoryID uuid.UUID
}

// ErrProductNotFound is returned when a product does not exist.
var ErrProductNotFound = errors.New("product not found")

// ErrProductConflict is returned when a product violates a unique
// constraint (e.g. a duplicate value in a unique column).
var ErrProductConflict = errors.New("product already exists")
