// Package domain holds the projects module's entities and value objects.
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Project is the project entity.
type Project struct {
	ID uuid.UUID
	Name string
	Key string
	Description string
	Archived bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreateProjectInput carries the data needed to create a project.
type CreateProjectInput struct {
	Name string
	Key string
	Description string
	Archived bool
}

// UpdateProjectInput carries the data needed to update a project.
type UpdateProjectInput struct {
	Name string
	Key string
	Description string
	Archived bool
}

// ErrProjectNotFound is returned when a project does not exist.
var ErrProjectNotFound = errors.New("project not found")

// ErrProjectConflict is returned when a project violates a unique
// constraint (e.g. a duplicate value in a unique column).
var ErrProjectConflict = errors.New("project already exists")
