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

type taskRepository struct {
	q *db.Queries
}

// NewTaskRepository constructs a TaskRepository.
func NewTaskRepository(q *db.Queries) ports.TaskRepository {
	return &taskRepository{q: q}
}

func (r *taskRepository) Create(ctx context.Context, in domain.CreateTaskInput) (domain.Task, error) {
	row, err := r.q.CreateTask(ctx, db.CreateTaskParams{
		Title: in.Title,
		Description: in.Description,
		Status: in.Status,
		Priority: in.Priority,
		DueDate: in.DueDate,
		Done: in.Done,
		ProjectID: in.ProjectID,
		AssigneeID: in.AssigneeID,
	})
	if err != nil {
		return domain.Task{}, mapTaskError(err)
	}
	return toTask(row), nil
}

func (r *taskRepository) Get(ctx context.Context, id uuid.UUID) (domain.Task, error) {
	row, err := r.q.GetTask(ctx, id)
	if err != nil {
		return domain.Task{}, mapTaskError(err)
	}
	return toTask(row), nil
}

func (r *taskRepository) List(ctx context.Context, limit, offset int32) ([]domain.Task, error) {
	rows, err := r.q.ListTasks(ctx, db.ListTasksParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	out := make([]domain.Task, 0, len(rows))
	for _, row := range rows {
		out = append(out, toTask(row))
	}
	return out, nil
}

func (r *taskRepository) Count(ctx context.Context) (int64, error) {
	return r.q.CountTasks(ctx)
}

func (r *taskRepository) Update(ctx context.Context, id uuid.UUID, in domain.UpdateTaskInput) (domain.Task, error) {
	row, err := r.q.UpdateTask(ctx, db.UpdateTaskParams{
		ID: id,
		Title: in.Title,
		Description: in.Description,
		Status: in.Status,
		Priority: in.Priority,
		DueDate: in.DueDate,
		Done: in.Done,
		ProjectID: in.ProjectID,
		AssigneeID: in.AssigneeID,
	})
	if err != nil {
		return domain.Task{}, mapTaskError(err)
	}
	return toTask(row), nil
}

func (r *taskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteTask(ctx, id)
}

func toTask(row db.Task) domain.Task {
	return domain.Task{
		ID: row.ID,
		Title: row.Title,
		Description: row.Description,
		Status: row.Status,
		Priority: row.Priority,
		DueDate: row.DueDate,
		Done: row.Done,
		ProjectID: row.ProjectID,
		AssigneeID: row.AssigneeID,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

// mapTaskError translates database errors into domain errors: a missing
// row becomes ErrTaskNotFound and a unique-constraint violation (SQLSTATE
// 23505) becomes ErrTaskConflict; anything else is returned unchanged.
func mapTaskError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrTaskNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return domain.ErrTaskConflict
	}
	return err
}
