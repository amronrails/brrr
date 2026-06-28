// Package bloghttp is the blog module's HTTP transport layer.
package bloghttp

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/example/blog/internal/modules/blog/internal/domain"
	"github.com/example/blog/internal/modules/blog/internal/services"
	"github.com/example/blog/internal/platform/auth"
	"github.com/example/blog/internal/platform/httpx"
)

type postRequest struct {
	Title string `json:"title" validate:"required"`
	Slug string `json:"slug" validate:"required"`
	Excerpt string `json:"excerpt"`
	Body string `json:"body"`
	Published bool `json:"published"`
	Views int32 `json:"views"`
	AuthorID uuid.UUID `json:"author_id" validate:"required"`
}

type postResponse struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Slug string `json:"slug"`
	Excerpt string `json:"excerpt"`
	Body string `json:"body"`
	Published bool `json:"published"`
	Views int32 `json:"views"`
	AuthorID uuid.UUID `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toPostResponse(m domain.Post) postResponse {
	return postResponse{
		ID: m.ID.String(),
		Title: m.Title,
		Slug: m.Slug,
		Excerpt: m.Excerpt,
		Body: m.Body,
		Published: m.Published,
		Views: m.Views,
		AuthorID: m.AuthorID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// PostHandler exposes CRUD endpoints for posts.
type PostHandler struct {
	svc      *services.PostService
	validate *validator.Validate
}

// NewPostHandler constructs a PostHandler.
func NewPostHandler(svc *services.PostService, validate *validator.Validate) *PostHandler {
	return &PostHandler{svc: svc, validate: validate}
}

func (h *PostHandler) create(w http.ResponseWriter, r *http.Request) {
	var req postRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		httpx.ValidationError(w, err)
		return
	}
	m, err := h.svc.Create(r.Context(), domain.CreatePostInput{
		Title: req.Title,
		Slug: req.Slug,
		Excerpt: req.Excerpt,
		Body: req.Body,
		Published: req.Published,
		Views: req.Views,
		AuthorID: req.AuthorID,
	})
	if err != nil {
		writePostError(w, err)
		return
	}
	httpx.JSON(w, http.StatusCreated, toPostResponse(m))
}

func (h *PostHandler) list(w http.ResponseWriter, r *http.Request) {
	limit := httpx.QueryInt(r, "limit", 20)
	offset := httpx.QueryInt(r, "offset", 0)
	items, total, err := h.svc.List(r.Context(), int32(limit), int32(offset))
	if err != nil {
		writePostError(w, err)
		return
	}
	out := make([]postResponse, 0, len(items))
	for _, m := range items {
		out = append(out, toPostResponse(m))
	}
	httpx.JSON(w, http.StatusOK, map[string]any{
		"data":   out,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *PostHandler) get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	m, err := h.svc.Get(r.Context(), id)
	if err != nil {
		writePostError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toPostResponse(m))
}

func (h *PostHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req postRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		httpx.ValidationError(w, err)
		return
	}
	m, err := h.svc.Update(r.Context(), id, domain.UpdatePostInput{
		Title: req.Title,
		Slug: req.Slug,
		Excerpt: req.Excerpt,
		Body: req.Body,
		Published: req.Published,
		Views: req.Views,
		AuthorID: req.AuthorID,
	})
	if err != nil {
		writePostError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toPostResponse(m))
}

func (h *PostHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		writePostError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// RegisterRoutes mounts the post CRUD endpoints (all require auth).
func (h *PostHandler) RegisterRoutes(r chi.Router, tokens *auth.TokenService) {
	r.Group(func(r chi.Router) {
		r.Use(tokens.RequireAuth)
		r.Post("/blog/posts", h.create)
		r.Get("/blog/posts", h.list)
		r.Get("/blog/posts/{id}", h.get)
		r.Put("/blog/posts/{id}", h.update)
		r.Delete("/blog/posts/{id}", h.delete)
	})
}

func writePostError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrPostNotFound):
		httpx.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrPostConflict):
		httpx.Error(w, http.StatusConflict, err.Error())
	default:
		httpx.Error(w, http.StatusInternalServerError, "internal server error")
	}
}
