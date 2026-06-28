// Package userhttp is the user module's HTTP transport layer: request/response
// DTOs, handlers and route registration.
package userhttp

import (
	"time"

	"github.com/example/shop/internal/modules/user/internal/domain"
	"github.com/example/shop/internal/modules/user/internal/services"
)

type registerRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
	Name     string `json:"name" validate:"max=100"`
}

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type userResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type authResponse struct {
	User         userResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	ExpiresAt    time.Time    `json:"expires_at"`
	RefreshToken string       `json:"refresh_token"`
}

func toUserResponse(u domain.User) userResponse {
	return userResponse{
		ID:        u.ID.String(),
		Email:     u.Email,
		Name:      u.Name,
		Role:      string(u.Role),
		CreatedAt: u.CreatedAt,
	}
}

func toAuthResponse(r services.AuthResult) authResponse {
	return authResponse{
		User:         toUserResponse(r.User),
		AccessToken:  r.AccessToken,
		ExpiresAt:    r.AccessExpiresAt,
		RefreshToken: r.RefreshToken,
	}
}
