// Package domain defines the user module's core entities and errors. It is the
// innermost layer and depends on nothing from the outer layers.
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Role enumerates the authorization roles a user may hold.
type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

// User is the core account entity.
type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	Name         string
	Role         Role
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// RefreshToken is a persisted (hashed) refresh token.
type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
}

// Domain errors are mapped to HTTP status codes by the transport layer.
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailTaken         = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrRefreshInvalid     = errors.New("invalid or expired refresh token")
)

// CreateUserParams are the inputs needed to persist a new user.
type CreateUserParams struct {
	Email        string
	PasswordHash string
	Name         string
	Role         Role
}

// StoreRefreshTokenParams are the inputs needed to persist a refresh token.
type StoreRefreshTokenParams struct {
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
}
