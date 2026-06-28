// Package postgres implements the projects module's ports on top of the
// centralized sqlc-generated query layer (internal/db).
package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/example/tasks/internal/db"
	"github.com/example/tasks/internal/modules/projects/internal/domain"
	"github.com/example/tasks/internal/modules/projects/internal/ports"
)

type projectRepository struct {
	q *db.Queries
}

// NewProjectRepository constructs a ProjectRepository.
func NewProjectRepository(q *db.Queries) ports.ProjectRepository {
	return &projectRepository{q: q}
}

func (r *projectRepository) Create(ctx context.Context, in domain.CreateProjectInput) (domain.Project, error) {
	row, err := r.q.CreateProject(ctx, db.CreateProjectParams{
		Name: in.Name,
		Key: in.Key,
		Description: in.Description,
		Archived: in.Archived,
	})
	if err != nil {
		return domain.Project{}, mapProjectError(err)
	}
	return toProject(row), nil
}

func (r *projectRepository) Get(ctx context.Context, id uuid.UUID) (domain.Project, error) {
	row, err := r.q.GetProject(ctx, id)
	if err != nil {
		return domain.Project{}, mapProjectError(err)
	}
	return toProject(row), nil
}

func (r *projectRepository) List(ctx context.Context, limit, offset int32) ([]domain.Project, error) {
	rows, err := r.q.ListProjects(ctx, db.ListProjectsParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	out := make([]domain.Project, 0, len(rows))
	for _, row := range rows {
		out = append(out, toProject(row))
	}
	return out, nil
}

func (r *projectRepository) Count(ctx context.Context) (int64, error) {
	return r.q.CountProjects(ctx)
}

func (r *projectRepository) Update(ctx context.Context, id uuid.UUID, in domain.UpdateProjectInput) (domain.Project, error) {
	row, err := r.q.UpdateProject(ctx, db.UpdateProjectParams{
		ID: id,
		Name: in.Name,
		Key: in.Key,
		Description: in.Description,
		Archived: in.Archived,
	})
	if err != nil {
		return domain.Project{}, mapProjectError(err)
	}
	return toProject(row), nil
}

func (r *projectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteProject(ctx, id)
}

func toProject(row db.Project) domain.Project {
	return domain.Project{
		ID: row.ID,
		Name: row.Name,
		Key: row.Key,
		Description: row.Description,
		Archived: row.Archived,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

// mapProjectError translates database errors into domain errors: a missing
// row becomes ErrProjectNotFound and a unique-constraint violation (SQLSTATE
// 23505) becomes ErrProjectConflict; anything else is returned unchanged.
func mapProjectError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrProjectNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return domain.ErrProjectConflict
	}
	return err
}
