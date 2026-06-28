// Package domain holds the projects module's entities and value objects.
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Label is the label entity.
type Label struct {
	ID uuid.UUID
	Name string
	Color string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreateLabelInput carries the data needed to create a label.
type CreateLabelInput struct {
	Name string
	Color string
}

// UpdateLabelInput carries the data needed to update a label.
type UpdateLabelInput struct {
	Name string
	Color string
}

// ErrLabelNotFound is returned when a label does not exist.
var ErrLabelNotFound = errors.New("label not found")

// ErrLabelConflict is returned when a label violates a unique
// constraint (e.g. a duplicate value in a unique column).
var ErrLabelConflict = errors.New("label already exists")
