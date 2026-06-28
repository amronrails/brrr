// Task public API DTO and accessor — the cross-module surface for
// tasks. Generated once by brrr; safe to extend by hand (the aggregate
// interface lives in api.go and is regenerated).
package projects

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/example/tasks/internal/modules/projects/internal/domain"
)

// Task is the public, cross-module representation of a task.
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

// TaskByID returns the task with the given id.
func (m *Module) TaskByID(ctx context.Context, id uuid.UUID) (Task, error) {
	e, err := m.taskSvc.Get(ctx, id)
	if err != nil {
		return Task{}, err
	}
	return toAPITask(e), nil
}

func toAPITask(e domain.Task) Task {
	return Task{
		ID: e.ID,
		Title: e.Title,
		Description: e.Description,
		Status: e.Status,
		Priority: e.Priority,
		DueDate: e.DueDate,
		Done: e.Done,
		ProjectID: e.ProjectID,
		AssigneeID: e.AssigneeID,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
