// Package projectshttp is the projects module's HTTP transport layer.
package projectshttp

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/example/tasks/internal/modules/projects/internal/domain"
	"github.com/example/tasks/internal/modules/projects/internal/services"
	"github.com/example/tasks/internal/platform/auth"
	"github.com/example/tasks/internal/platform/httpx"
)

type labelRequest struct {
	Name string `json:"name" validate:"required"`
	Color string `json:"color"`
}

type labelResponse struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Color string `json:"color"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toLabelResponse(m domain.Label) labelResponse {
	return labelResponse{
		ID: m.ID.String(),
		Name: m.Name,
		Color: m.Color,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// LabelHandler exposes CRUD endpoints for labels.
type LabelHandler struct {
	svc      *services.LabelService
	validate *validator.Validate
}

// NewLabelHandler constructs a LabelHandler.
func NewLabelHandler(svc *services.LabelService, validate *validator.Validate) *LabelHandler {
	return &LabelHandler{svc: svc, validate: validate}
}

func (h *LabelHandler) create(w http.ResponseWriter, r *http.Request) {
	var req labelRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		httpx.ValidationError(w, err)
		return
	}
	m, err := h.svc.Create(r.Context(), domain.CreateLabelInput{
		Name: req.Name,
		Color: req.Color,
	})
	if err != nil {
		writeLabelError(w, err)
		return
	}
	httpx.JSON(w, http.StatusCreated, toLabelResponse(m))
}

func (h *LabelHandler) list(w http.ResponseWriter, r *http.Request) {
	limit := httpx.QueryInt(r, "limit", 20)
	offset := httpx.QueryInt(r, "offset", 0)
	items, total, err := h.svc.List(r.Context(), int32(limit), int32(offset))
	if err != nil {
		writeLabelError(w, err)
		return
	}
	out := make([]labelResponse, 0, len(items))
	for _, m := range items {
		out = append(out, toLabelResponse(m))
	}
	httpx.JSON(w, http.StatusOK, map[string]any{
		"data":   out,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *LabelHandler) get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	m, err := h.svc.Get(r.Context(), id)
	if err != nil {
		writeLabelError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toLabelResponse(m))
}

func (h *LabelHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req labelRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		httpx.ValidationError(w, err)
		return
	}
	m, err := h.svc.Update(r.Context(), id, domain.UpdateLabelInput{
		Name: req.Name,
		Color: req.Color,
	})
	if err != nil {
		writeLabelError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toLabelResponse(m))
}

func (h *LabelHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeLabelError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// RegisterRoutes mounts the label CRUD endpoints (all require auth).
func (h *LabelHandler) RegisterRoutes(r chi.Router, tokens *auth.TokenService) {
	r.Group(func(r chi.Router) {
		r.Use(tokens.RequireAuth)
		r.Post("/projects/labels", h.create)
		r.Get("/projects/labels", h.list)
		r.Get("/projects/labels/{id}", h.get)
		r.Put("/projects/labels/{id}", h.update)
		r.Delete("/projects/labels/{id}", h.delete)
	})
}

func writeLabelError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrLabelNotFound):
		httpx.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrLabelConflict):
		httpx.Error(w, http.StatusConflict, err.Error())
	default:
		httpx.Error(w, http.StatusInternalServerError, "internal server error")
	}
}
