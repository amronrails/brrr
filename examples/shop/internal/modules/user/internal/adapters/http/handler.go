package userhttp

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"

	"github.com/example/shop/internal/modules/user/internal/domain"
	"github.com/example/shop/internal/modules/user/internal/services"
	"github.com/example/shop/internal/platform/auth"
	"github.com/example/shop/internal/platform/httpx"
)

// Handler exposes the user module's HTTP endpoints.
type Handler struct {
	svc      *services.Service
	validate *validator.Validate
}

// NewHandler constructs a Handler.
func NewHandler(svc *services.Service, validate *validator.Validate) *Handler {
	return &Handler{svc: svc, validate: validate}
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		httpx.ValidationError(w, err)
		return
	}
	res, err := h.svc.Register(r.Context(), services.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.JSON(w, http.StatusCreated, toAuthResponse(res))
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		httpx.ValidationError(w, err)
		return
	}
	res, err := h.svc.Login(r.Context(), services.LoginInput{Email: req.Email, Password: req.Password})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toAuthResponse(res))
}

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		httpx.ValidationError(w, err)
		return
	}
	res, err := h.svc.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toAuthResponse(res))
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.svc.Logout(r.Context(), req.RefreshToken); err != nil {
		writeServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) me(w http.ResponseWriter, r *http.Request) {
	id, ok := auth.IdentityFrom(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "authentication required")
		return
	}
	u, err := h.svc.Me(r.Context(), id.UserID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toUserResponse(u))
}

func (h *Handler) listUsers(w http.ResponseWriter, r *http.Request) {
	limit := atoiDefault(r.URL.Query().Get("limit"), 20)
	offset := atoiDefault(r.URL.Query().Get("offset"), 0)
	users, err := h.svc.ListUsers(r.Context(), int32(limit), int32(offset))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	out := make([]userResponse, 0, len(users))
	for _, u := range users {
		out = append(out, toUserResponse(u))
	}
	httpx.JSON(w, http.StatusOK, map[string]any{"data": out})
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrEmailTaken):
		httpx.Error(w, http.StatusConflict, err.Error())
	case errors.Is(err, domain.ErrInvalidCredentials), errors.Is(err, domain.ErrRefreshInvalid):
		httpx.Error(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, domain.ErrUserNotFound):
		httpx.Error(w, http.StatusNotFound, err.Error())
	default:
		httpx.Error(w, http.StatusInternalServerError, "internal server error")
	}
}

func atoiDefault(s string, def int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}
