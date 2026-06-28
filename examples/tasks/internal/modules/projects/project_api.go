// Project public API DTO and accessor — the cross-module surface for
// projects. Generated once by brrr; safe to extend by hand (the aggregate
// interface lives in api.go and is regenerated).
package projects

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/example/tasks/internal/modules/projects/internal/domain"
)

// Project is the public, cross-module representation of a project.
type Project struct {
	ID uuid.UUID
	Name string
	Key string
	Description string
	Archived bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ProjectByID returns the project with the given id.
func (m *Module) ProjectByID(ctx context.Context, id uuid.UUID) (Project, error) {
	e, err := m.projectSvc.Get(ctx, id)
	if err != nil {
		return Project{}, err
	}
	return toAPIProject(e), nil
}

func toAPIProject(e domain.Project) Project {
	return Project{
		ID: e.ID,
		Name: e.Name,
		Key: e.Key,
		Description: e.Description,
		Archived: e.Archived,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
