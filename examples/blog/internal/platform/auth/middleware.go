package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/example/blog/internal/platform/httpx"
)

type ctxKey int

const identityKey ctxKey = iota

// Identity is the authenticated principal extracted from an access token.
type Identity struct {
	UserID uuid.UUID
	Role   string
}

// RequireAuth is middleware that validates the bearer access token and stores
// the resulting Identity in the request context.
func (s *TokenService) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := bearerToken(r)
		if raw == "" {
			httpx.Error(w, http.StatusUnauthorized, "missing bearer token")
			return
		}
		claims, err := s.ParseAccess(raw)
		if err != nil {
			httpx.Error(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}
		uid, err := uuid.Parse(claims.Subject)
		if err != nil {
			httpx.Error(w, http.StatusUnauthorized, "invalid token subject")
			return
		}
		ctx := context.WithValue(r.Context(), identityKey, Identity{UserID: uid, Role: claims.Role})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole returns middleware that allows the request only if the
// authenticated identity has the given role. It must be chained after
// RequireAuth.
func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id, ok := IdentityFrom(r.Context())
			if !ok {
				httpx.Error(w, http.StatusUnauthorized, "authentication required")
				return
			}
			if id.Role != role {
				httpx.Error(w, http.StatusForbidden, "insufficient permissions")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// IdentityFrom returns the authenticated Identity stored in ctx, if any.
func IdentityFrom(ctx context.Context) (Identity, bool) {
	id, ok := ctx.Value(identityKey).(Identity)
	return id, ok
}

func bearerToken(r *http.Request) string {
	h := r.Header.Get("Authorization")
	if h == "" {
		return ""
	}
	parts := strings.SplitN(h, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
