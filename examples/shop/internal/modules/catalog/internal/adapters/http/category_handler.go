// Package cataloghttp is the catalog module's HTTP transport layer.
package cataloghttp

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/example/shop/internal/modules/catalog/internal/domain"
	"github.com/example/shop/internal/modules/catalog/internal/services"
	"github.com/example/shop/internal/platform/auth"
	"github.com/example/shop/internal/platform/httpx"
)

type categoryRequest struct {
	Name string `json:"name" validate:"required"`
	Slug string `json:"slug" validate:"required"`
}

type categoryResponse struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toCategoryResponse(m domain.Category) categoryResponse {
	return categoryResponse{
		ID: m.ID.String(),
		Name: m.Name,
		Slug: m.Slug,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// CategoryHandler exposes CRUD endpoints for categories.
type CategoryHandler struct {
	svc      *services.CategoryService
	validate *validator.Validate
}

// NewCategoryHandler constructs a CategoryHandler.
func NewCategoryHandler(svc *services.CategoryService, validate *validator.Validate) *CategoryHandler {
	return &CategoryHandler{svc: svc, validate: validate}
}

func (h *CategoryHandler) create(w http.ResponseWriter, r *http.Request) {
	var req categoryRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		httpx.ValidationError(w, err)
		return
	}
	m, err := h.svc.Create(r.Context(), domain.CreateCategoryInput{
		Name: req.Name,
		Slug: req.Slug,
	})
	if err != nil {
		writeCategoryError(w, err)
		return
	}
	httpx.JSON(w, http.StatusCreated, toCategoryResponse(m))
}

func (h *CategoryHandler) list(w http.ResponseWriter, r *http.Request) {
	limit := httpx.QueryInt(r, "limit", 20)
	offset := httpx.QueryInt(r, "offset", 0)
	items, total, err := h.svc.List(r.Context(), int32(limit), int32(offset))
	if err != nil {
		writeCategoryError(w, err)
		return
	}
	out := make([]categoryResponse, 0, len(items))
	for _, m := range items {
		out = append(out, toCategoryResponse(m))
	}
	httpx.JSON(w, http.StatusOK, map[string]any{
		"data":   out,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *CategoryHandler) get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	m, err := h.svc.Get(r.Context(), id)
	if err != nil {
		writeCategoryError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toCategoryResponse(m))
}

func (h *CategoryHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req categoryRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		httpx.ValidationError(w, err)
		return
	}
	m, err := h.svc.Update(r.Context(), id, domain.UpdateCategoryInput{
		Name: req.Name,
		Slug: req.Slug,
	})
	if err != nil {
		writeCategoryError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toCategoryResponse(m))
}

func (h *CategoryHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeCategoryError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// RegisterRoutes mounts the category CRUD endpoints (all require auth).
func (h *CategoryHandler) RegisterRoutes(r chi.Router, tokens *auth.TokenService) {
	r.Group(func(r chi.Router) {
		r.Use(tokens.RequireAuth)
		r.Post("/catalog/categories", h.create)
		r.Get("/catalog/categories", h.list)
		r.Get("/catalog/categories/{id}", h.get)
		r.Put("/catalog/categories/{id}", h.update)
		r.Delete("/catalog/categories/{id}", h.delete)
	})
}

func writeCategoryError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrCategoryNotFound):
		httpx.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrCategoryConflict):
		httpx.Error(w, http.StatusConflict, err.Error())
	default:
		httpx.Error(w, http.StatusInternalServerError, "internal server error")
	}
}
