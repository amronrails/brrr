// Package services holds the user module's business logic: registration, login,
// token rotation and admin queries. It depends only on the domain entities, the
// ports it consumes, and platform auth primitives.
package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/example/tasks/internal/modules/user/internal/domain"
	"github.com/example/tasks/internal/modules/user/internal/ports"
	"github.com/example/tasks/internal/platform/auth"
)

// Service implements the user module's use cases.
type Service struct {
	repo   ports.UserRepository
	tokens *auth.TokenService
}

// New constructs a Service.
func New(repo ports.UserRepository, tokens *auth.TokenService) *Service {
	return &Service{repo: repo, tokens: tokens}
}

// RegisterInput carries a registration request.
type RegisterInput struct {
	Email    string
	Password string
	Name     string
}

// LoginInput carries a login request.
type LoginInput struct {
	Email    string
	Password string
}

// AuthResult is returned by the authentication use cases.
type AuthResult struct {
	User            domain.User
	AccessToken     string
	AccessExpiresAt time.Time
	RefreshToken    string
}

// Register creates a new account. The very first account in the system is
// promoted to admin so the dashboard is reachable after a fresh install.
func (s *Service) Register(ctx context.Context, in RegisterInput) (AuthResult, error) {
	hash, err := auth.HashPassword(in.Password)
	if err != nil {
		return AuthResult{}, err
	}

	role := domain.RoleUser
	if count, err := s.repo.CountUsers(ctx); err != nil {
		return AuthResult{}, err
	} else if count == 0 {
		role = domain.RoleAdmin
	}

	u, err := s.repo.CreateUser(ctx, domain.CreateUserParams{
		Email:        normalizeEmail(in.Email),
		PasswordHash: hash,
		Name:         strings.TrimSpace(in.Name),
		Role:         role,
	})
	if err != nil {
		return AuthResult{}, err
	}
	return s.issueTokens(ctx, u)
}

// Login authenticates by email and password.
func (s *Service) Login(ctx context.Context, in LoginInput) (AuthResult, error) {
	u, err := s.repo.GetByEmail(ctx, normalizeEmail(in.Email))
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return AuthResult{}, domain.ErrInvalidCredentials
		}
		return AuthResult{}, err
	}
	if !auth.CheckPassword(u.PasswordHash, in.Password) {
		return AuthResult{}, domain.ErrInvalidCredentials
	}
	return s.issueTokens(ctx, u)
}

// Refresh validates a refresh token, rotates it, and issues a new token pair.
func (s *Service) Refresh(ctx context.Context, refreshToken string) (AuthResult, error) {
	hash := auth.HashToken(refreshToken)
	rt, err := s.repo.GetRefreshToken(ctx, hash)
	if err != nil {
		return AuthResult{}, err
	}
	if time.Now().After(rt.ExpiresAt) {
		_ = s.repo.DeleteRefreshToken(ctx, hash)
		return AuthResult{}, domain.ErrRefreshInvalid
	}
	u, err := s.repo.GetByID(ctx, rt.UserID)
	if err != nil {
		return AuthResult{}, err
	}
	if err := s.repo.DeleteRefreshToken(ctx, hash); err != nil {
		return AuthResult{}, err
	}
	return s.issueTokens(ctx, u)
}

// Logout revokes a single refresh token.
func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	return s.repo.DeleteRefreshToken(ctx, auth.HashToken(refreshToken))
}

// Me returns the authenticated user's profile.
func (s *Service) Me(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	return s.repo.GetByID(ctx, userID)
}

// ListUsers returns a page of users (admin only at the transport layer).
func (s *Service) ListUsers(ctx context.Context, limit, offset int32) ([]domain.User, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.ListUsers(ctx, limit, offset)
}

func (s *Service) issueTokens(ctx context.Context, u domain.User) (AuthResult, error) {
	access, exp, err := s.tokens.GenerateAccess(u.ID, string(u.Role))
	if err != nil {
		return AuthResult{}, err
	}
	refresh, hash, rexp, err := s.tokens.GenerateRefresh()
	if err != nil {
		return AuthResult{}, err
	}
	if err := s.repo.StoreRefreshToken(ctx, domain.StoreRefreshTokenParams{
		UserID:    u.ID,
		TokenHash: hash,
		ExpiresAt: rexp,
	}); err != nil {
		return AuthResult{}, err
	}
	return AuthResult{User: u, AccessToken: access, AccessExpiresAt: exp, RefreshToken: refresh}, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
