// Package ports declares the interfaces the user module's services depend on.
// They are consumer-defined contracts: the services own them and the adapters
// (e.g. postgres) implement them, keeping the dependency arrow pointing inward.
package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/example/blog/internal/modules/user/internal/domain"
)

// UserRepository is the persistence contract for the user module. It is
// implemented by the postgres adapter and consumed by the service layer.
type UserRepository interface {
	CreateUser(ctx context.Context, p domain.CreateUserParams) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	CountUsers(ctx context.Context) (int64, error)
	ListUsers(ctx context.Context, limit, offset int32) ([]domain.User, error)

	StoreRefreshToken(ctx context.Context, p domain.StoreRefreshTokenParams) error
	GetRefreshToken(ctx context.Context, hash string) (domain.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, hash string) error
	DeleteUserRefreshTokens(ctx context.Context, userID uuid.UUID) error
}
