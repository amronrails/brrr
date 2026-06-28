// Package domain holds the catalog module's entities and value objects.
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Category is the category entity.
type Category struct {
	ID uuid.UUID
	Name string
	Slug string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreateCategoryInput carries the data needed to create a category.
type CreateCategoryInput struct {
	Name string
	Slug string
}

// UpdateCategoryInput carries the data needed to update a category.
type UpdateCategoryInput struct {
	Name string
	Slug string
}

// ErrCategoryNotFound is returned when a category does not exist.
var ErrCategoryNotFound = errors.New("category not found")

// ErrCategoryConflict is returned when a category violates a unique
// constraint (e.g. a duplicate value in a unique column).
var ErrCategoryConflict = errors.New("category already exists")
