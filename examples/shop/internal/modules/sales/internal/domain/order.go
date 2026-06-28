// Package domain holds the sales module's entities and value objects.
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Order is the order entity.
type Order struct {
	ID uuid.UUID
	Status string
	Total string
	PlacedAt time.Time
	CustomerID uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreateOrderInput carries the data needed to create a order.
type CreateOrderInput struct {
	Status string
	Total string
	PlacedAt time.Time
	CustomerID uuid.UUID
}

// UpdateOrderInput carries the data needed to update a order.
type UpdateOrderInput struct {
	Status string
	Total string
	PlacedAt time.Time
	CustomerID uuid.UUID
}

// ErrOrderNotFound is returned when a order does not exist.
var ErrOrderNotFound = errors.New("order not found")

// ErrOrderConflict is returned when a order violates a unique
// constraint (e.g. a duplicate value in a unique column).
var ErrOrderConflict = errors.New("order already exists")
