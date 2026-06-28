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

type taskRequest struct {
	Title string `json:"title" validate:"required"`
	Description string `json:"description"`
	Status string `json:"status" validate:"required"`
	Priority int32 `json:"priority"`
	DueDate time.Time `json:"due_date"`
	Done bool `json:"done"`
	ProjectID uuid.UUID `json:"project_id" validate:"required"`
	AssigneeID uuid.UUID `json:"assignee_id" validate:"required"`
}

type taskResponse struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Status string `json:"status"`
	Priority int32 `json:"priority"`
	DueDate time.Time `json:"due_date"`
	Done bool `json:"done"`
	ProjectID uuid.UUID `json:"project_id"`
	AssigneeID uuid.UUID `json:"assignee_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toTaskResponse(m domain.Task) taskResponse {
	return taskResponse{
		ID: m.ID.String(),
		Title: m.Title,
		Description: m.Description,
		Status: m.Status,
		Priority: m.Priority,
		DueDate: m.DueDate,
		Done: m.Done,
		ProjectID: m.ProjectID,
		AssigneeID: m.AssigneeID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// TaskHandler exposes CRUD endpoints for tasks.
type TaskHandler struct {
	svc      *services.TaskService
	validate *validator.Validate
}

// NewTaskHandler constructs a TaskHandler.
func NewTaskHandler(svc *services.TaskService, validate *validator.Validate) *TaskHandler {
	return &TaskHandler{svc: svc, validate: validate}
}

func (h *TaskHandler) create(w http.ResponseWriter, r *http.Request) {
	var req taskRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		httpx.ValidationError(w, err)
		return
	}
	m, err := h.svc.Create(r.Context(), domain.CreateTaskInput{
		Title: req.Title,
		Description: req.Description,
		Status: req.Status,
		Priority: req.Priority,
		DueDate: req.DueDate,
		Done: req.Done,
		ProjectID: req.ProjectID,
		AssigneeID: req.AssigneeID,
	})
	if err != nil {
		writeTaskError(w, err)
		return
	}
	httpx.JSON(w, http.StatusCreated, toTaskResponse(m))
}

func (h *TaskHandler) list(w http.ResponseWriter, r *http.Request) {
	limit := httpx.QueryInt(r, "limit", 20)
	offset := httpx.QueryInt(r, "offset", 0)
	items, total, err := h.svc.List(r.Context(), int32(limit), int32(offset))
	if err != nil {
		writeTaskError(w, err)
		return
	}
	out := make([]taskResponse, 0, len(items))
	for _, m := range items {
		out = append(out, toTaskResponse(m))
	}
	httpx.JSON(w, http.StatusOK, map[string]any{
		"data":   out,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *TaskHandler) get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	m, err := h.svc.Get(r.Context(), id)
	if err != nil {
		writeTaskError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toTaskResponse(m))
}

func (h *TaskHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req taskRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		httpx.ValidationError(w, err)
		return
	}
	m, err := h.svc.Update(r.Context(), id, domain.UpdateTaskInput{
		Title: req.Title,
		Description: req.Description,
		Status: req.Status,
		Priority: req.Priority,
		DueDate: req.DueDate,
		Done: req.Done,
		ProjectID: req.ProjectID,
		AssigneeID: req.AssigneeID,
	})
	if err != nil {
		writeTaskError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toTaskResponse(m))
}

func (h *TaskHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeTaskError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// RegisterRoutes mounts the task CRUD endpoints (all require auth).
func (h *TaskHandler) RegisterRoutes(r chi.Router, tokens *auth.TokenService) {
	r.Group(func(r chi.Router) {
		r.Use(tokens.RequireAuth)
		r.Post("/projects/tasks", h.create)
		r.Get("/projects/tasks", h.list)
		r.Get("/projects/tasks/{id}", h.get)
		r.Put("/projects/tasks/{id}", h.update)
		r.Delete("/projects/tasks/{id}", h.delete)
	})
}

func writeTaskError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrTaskNotFound):
		httpx.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrTaskConflict):
		httpx.Error(w, http.StatusConflict, err.Error())
	default:
		httpx.Error(w, http.StatusInternalServerError, "internal server error")
	}
}
