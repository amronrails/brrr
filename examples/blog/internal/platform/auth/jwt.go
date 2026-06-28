// Package auth provides password hashing, JWT access tokens, opaque refresh
// tokens, and HTTP middleware for authentication and authorization.
package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/example/blog/internal/platform/config"
)

// AccessClaims are the custom claims carried by an access token.
type AccessClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// TokenService issues and verifies access and refresh tokens.
//
// Access tokens are short-lived signed JWTs (HS256). Refresh tokens are opaque
// random strings; only their SHA-256 hash is persisted, so they can be revoked
// and rotated server-side.
type TokenService struct {
	secret     []byte
	issuer     string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

// NewTokenService builds a TokenService from JWT configuration.
func NewTokenService(cfg config.JWTConfig) *TokenService {
	return &TokenService{
		secret:     []byte(cfg.Secret),
		issuer:     cfg.Issuer,
		accessTTL:  cfg.AccessTTL,
		refreshTTL: cfg.RefreshTTL,
	}
}

// GenerateAccess issues a signed access token for a user, returning the token
// and its expiry.
func (s *TokenService) GenerateAccess(userID uuid.UUID, role string) (string, time.Time, error) {
	now := time.Now()
	exp := now.Add(s.accessTTL)
	claims := AccessClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			Issuer:    s.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secret)
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, exp, nil
}

// ParseAccess validates an access token and returns its claims.
func (s *TokenService) ParseAccess(token string) (*AccessClaims, error) {
	claims := &AccessClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.secret, nil
	}, jwt.WithIssuer(s.issuer), jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		return nil, err
	}
	return claims, nil
}

// GenerateRefresh creates a new opaque refresh token, returning the plaintext
// token (handed to the client), its hash (persisted), and its expiry.
func (s *TokenService) GenerateRefresh() (token, hash string, exp time.Time, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return "", "", time.Time{}, err
	}
	token = base64.RawURLEncoding.EncodeToString(b)
	return token, HashToken(token), time.Now().Add(s.refreshTTL), nil
}

// HashToken returns the hex-encoded SHA-256 hash of a refresh token, used as the
// stored lookup key.
func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
