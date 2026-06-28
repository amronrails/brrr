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

type labelRepository struct {
	q *db.Queries
}

// NewLabelRepository constructs a LabelRepository.
func NewLabelRepository(q *db.Queries) ports.LabelRepository {
	return &labelRepository{q: q}
}

func (r *labelRepository) Create(ctx context.Context, in domain.CreateLabelInput) (domain.Label, error) {
	row, err := r.q.CreateLabel(ctx, db.CreateLabelParams{
		Name: in.Name,
		Color: in.Color,
	})
	if err != nil {
		return domain.Label{}, mapLabelError(err)
	}
	return toLabel(row), nil
}

func (r *labelRepository) Get(ctx context.Context, id uuid.UUID) (domain.Label, error) {
	row, err := r.q.GetLabel(ctx, id)
	if err != nil {
		return domain.Label{}, mapLabelError(err)
	}
	return toLabel(row), nil
}

func (r *labelRepository) List(ctx context.Context, limit, offset int32) ([]domain.Label, error) {
	rows, err := r.q.ListLabels(ctx, db.ListLabelsParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	out := make([]domain.Label, 0, len(rows))
	for _, row := range rows {
		out = append(out, toLabel(row))
	}
	return out, nil
}

func (r *labelRepository) Count(ctx context.Context) (int64, error) {
	return r.q.CountLabels(ctx)
}

func (r *labelRepository) Update(ctx context.Context, id uuid.UUID, in domain.UpdateLabelInput) (domain.Label, error) {
	row, err := r.q.UpdateLabel(ctx, db.UpdateLabelParams{
		ID: id,
		Name: in.Name,
		Color: in.Color,
	})
	if err != nil {
		return domain.Label{}, mapLabelError(err)
	}
	return toLabel(row), nil
}

func (r *labelRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteLabel(ctx, id)
}

func toLabel(row db.Label) domain.Label {
	return domain.Label{
		ID: row.ID,
		Name: row.Name,
		Color: row.Color,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

// mapLabelError translates database errors into domain errors: a missing
// row becomes ErrLabelNotFound and a unique-constraint violation (SQLSTATE
// 23505) becomes ErrLabelConflict; anything else is returned unchanged.
func mapLabelError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrLabelNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return domain.ErrLabelConflict
	}
	return err
}
