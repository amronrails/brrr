// Package domain holds the sales module's entities and value objects.
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// OrderItem is the order_item entity.
type OrderItem struct {
	ID uuid.UUID
	Quantity int32
	UnitPrice string
	OrderID uuid.UUID
	ProductID uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreateOrderItemInput carries the data needed to create a order_item.
type CreateOrderItemInput struct {
	Quantity int32
	UnitPrice string
	OrderID uuid.UUID
	ProductID uuid.UUID
}

// UpdateOrderItemInput carries the data needed to update a order_item.
type UpdateOrderItemInput struct {
	Quantity int32
	UnitPrice string
	OrderID uuid.UUID
	ProductID uuid.UUID
}

// ErrOrderItemNotFound is returned when a order_item does not exist.
var ErrOrderItemNotFound = errors.New("order_item not found")

// ErrOrderItemConflict is returned when a order_item violates a unique
// constraint (e.g. a duplicate value in a unique column).
var ErrOrderItemConflict = errors.New("order_item already exists")
