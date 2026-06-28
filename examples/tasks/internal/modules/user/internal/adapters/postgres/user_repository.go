// Package postgres implements the user module's ports on top of the centralized
// sqlc-generated query layer (internal/db).
package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/example/tasks/internal/db"
	"github.com/example/tasks/internal/modules/user/internal/domain"
	"github.com/example/tasks/internal/modules/user/internal/ports"
)

// Repository adapts the generated *db.Queries to the ports.UserRepository
// interface, translating database rows and errors into domain types.
type Repository struct {
	q *db.Queries
}

// New constructs a Repository backed by the given sqlc queries.
func New(q *db.Queries) *Repository { return &Repository{q: q} }

var _ ports.UserRepository = (*Repository)(nil)

func (r *Repository) CreateUser(ctx context.Context, p domain.CreateUserParams) (domain.User, error) {
	row, err := r.q.CreateUser(ctx, db.CreateUserParams{
		Email:        p.Email,
		PasswordHash: p.PasswordHash,
		Name:         p.Name,
		Role:         string(p.Role),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.User{}, domain.ErrEmailTaken
		}
		return domain.User{}, err
	}
	return toUser(row), nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	row, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}
	return toUser(row), nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	row, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}
	return toUser(row), nil
}

func (r *Repository) CountUsers(ctx context.Context) (int64, error) {
	return r.q.CountUsers(ctx)
}

func (r *Repository) ListUsers(ctx context.Context, limit, offset int32) ([]domain.User, error) {
	rows, err := r.q.ListUsers(ctx, db.ListUsersParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	users := make([]domain.User, 0, len(rows))
	for _, row := range rows {
		users = append(users, toUser(row))
	}
	return users, nil
}

func (r *Repository) StoreRefreshToken(ctx context.Context, p domain.StoreRefreshTokenParams) error {
	return r.q.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
		UserID:    p.UserID,
		TokenHash: p.TokenHash,
		ExpiresAt: p.ExpiresAt,
	})
}

func (r *Repository) GetRefreshToken(ctx context.Context, hash string) (domain.RefreshToken, error) {
	row, err := r.q.GetRefreshToken(ctx, hash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.RefreshToken{}, domain.ErrRefreshInvalid
		}
		return domain.RefreshToken{}, err
	}
	return domain.RefreshToken{
		ID:        row.ID,
		UserID:    row.UserID,
		TokenHash: row.TokenHash,
		ExpiresAt: row.ExpiresAt,
		CreatedAt: row.CreatedAt,
	}, nil
}

func (r *Repository) DeleteRefreshToken(ctx context.Context, hash string) error {
	return r.q.DeleteRefreshToken(ctx, hash)
}

func (r *Repository) DeleteUserRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	return r.q.DeleteUserRefreshTokens(ctx, userID)
}

func toUser(row db.User) domain.User {
	return domain.User{
		ID:           row.ID,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		Name:         row.Name,
		Role:         domain.Role(row.Role),
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
}
