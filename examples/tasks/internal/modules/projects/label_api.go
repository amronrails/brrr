// Label public API DTO and accessor — the cross-module surface for
// labels. Generated once by brrr; safe to extend by hand (the aggregate
// interface lives in api.go and is regenerated).
package projects

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/example/tasks/internal/modules/projects/internal/domain"
)

// Label is the public, cross-module representation of a label.
type Label struct {
	ID uuid.UUID
	Name string
	Color string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// LabelByID returns the label with the given id.
func (m *Module) LabelByID(ctx context.Context, id uuid.UUID) (Label, error) {
	e, err := m.labelSvc.Get(ctx, id)
	if err != nil {
		return Label{}, err
	}
	return toAPILabel(e), nil
}

func toAPILabel(e domain.Label) Label {
	return Label{
		ID: e.ID,
		Name: e.Name,
		Color: e.Color,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
