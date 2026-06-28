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

type projectRequest struct {
	Name string `json:"name" validate:"required"`
	Key string `json:"key" validate:"required"`
	Description string `json:"description"`
	Archived bool `json:"archived"`
}

type projectResponse struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Key string `json:"key"`
	Description string `json:"description"`
	Archived bool `json:"archived"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toProjectResponse(m domain.Project) projectResponse {
	return projectResponse{
		ID: m.ID.String(),
		Name: m.Name,
		Key: m.Key,
		Description: m.Description,
		Archived: m.Archived,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// ProjectHandler exposes CRUD endpoints for projects.
type ProjectHandler struct {
	svc      *services.ProjectService
	validate *validator.Validate
}

// NewProjectHandler constructs a ProjectHandler.
func NewProjectHandler(svc *services.ProjectService, validate *validator.Validate) *ProjectHandler {
	return &ProjectHandler{svc: svc, validate: validate}
}

func (h *ProjectHandler) create(w http.ResponseWriter, r *http.Request) {
	var req projectRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		httpx.ValidationError(w, err)
		return
	}
	m, err := h.svc.Create(r.Context(), domain.CreateProjectInput{
		Name: req.Name,
		Key: req.Key,
		Description: req.Description,
		Archived: req.Archived,
	})
	if err != nil {
		writeProjectError(w, err)
		return
	}
	httpx.JSON(w, http.StatusCreated, toProjectResponse(m))
}

func (h *ProjectHandler) list(w http.ResponseWriter, r *http.Request) {
	limit := httpx.QueryInt(r, "limit", 20)
	offset := httpx.QueryInt(r, "offset", 0)
	items, total, err := h.svc.List(r.Context(), int32(limit), int32(offset))
	if err != nil {
		writeProjectError(w, err)
		return
	}
	out := make([]projectResponse, 0, len(items))
	for _, m := range items {
		out = append(out, toProjectResponse(m))
	}
	httpx.JSON(w, http.StatusOK, map[string]any{
		"data":   out,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *ProjectHandler) get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	m, err := h.svc.Get(r.Context(), id)
	if err != nil {
		writeProjectError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toProjectResponse(m))
}

func (h *ProjectHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req projectRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		httpx.ValidationError(w, err)
		return
	}
	m, err := h.svc.Update(r.Context(), id, domain.UpdateProjectInput{
		Name: req.Name,
		Key: req.Key,
		Description: req.Description,
		Archived: req.Archived,
	})
	if err != nil {
		writeProjectError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toProjectResponse(m))
}

func (h *ProjectHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeProjectError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// RegisterRoutes mounts the project CRUD endpoints (all require auth).
func (h *ProjectHandler) RegisterRoutes(r chi.Router, tokens *auth.TokenService) {
	r.Group(func(r chi.Router) {
		r.Use(tokens.RequireAuth)
		r.Post("/projects/projects", h.create)
		r.Get("/projects/projects", h.list)
		r.Get("/projects/projects/{id}", h.get)
		r.Put("/projects/projects/{id}", h.update)
		r.Delete("/projects/projects/{id}", h.delete)
	})
}

func writeProjectError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrProjectNotFound):
		httpx.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrProjectConflict):
		httpx.Error(w, http.StatusConflict, err.Error())
	default:
		httpx.Error(w, http.StatusInternalServerError, "internal server error")
	}
}
