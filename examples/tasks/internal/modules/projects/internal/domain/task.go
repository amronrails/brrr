// Package domain holds the projects module's entities and value objects.
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Task is the task entity.
type Task struct {
	ID uuid.UUID
	Title string
	Description string
	Status string
	Priority int32
	DueDate time.Time
	Done bool
	ProjectID uuid.UUID
	AssigneeID uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreateTaskInput carries the data needed to create a task.
type CreateTaskInput struct {
	Title string
	Description string
	Status string
	Priority int32
	DueDate time.Time
	Done bool
	ProjectID uuid.UUID
	AssigneeID uuid.UUID
}

// UpdateTaskInput carries the data needed to update a task.
type UpdateTaskInput struct {
	Title string
	Description string
	Status string
	Priority int32
	DueDate time.Time
	Done bool
	ProjectID uuid.UUID
	AssigneeID uuid.UUID
}

// ErrTaskNotFound is returned when a task does not exist.
var ErrTaskNotFound = errors.New("task not found")

// ErrTaskConflict is returned when a task violates a unique
// constraint (e.g. a duplicate value in a unique column).
var ErrTaskConflict = errors.New("task already exists")
