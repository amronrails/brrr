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

type commentRequest struct {
	Body string `json:"body" validate:"required"`
	Approved bool `json:"approved"`
	PostID uuid.UUID `json:"post_id" validate:"required"`
	AuthorID uuid.UUID `json:"author_id" validate:"required"`
}

type commentResponse struct {
	ID string `json:"id"`
	Body string `json:"body"`
	Approved bool `json:"approved"`
	PostID uuid.UUID `json:"post_id"`
	AuthorID uuid.UUID `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toCommentResponse(m domain.Comment) commentResponse {
	return commentResponse{
		ID: m.ID.String(),
		Body: m.Body,
		Approved: m.Approved,
		PostID: m.PostID,
		AuthorID: m.AuthorID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// CommentHandler exposes CRUD endpoints for comments.
type CommentHandler struct {
	svc      *services.CommentService
	validate *validator.Validate
}

// NewCommentHandler constructs a CommentHandler.
func NewCommentHandler(svc *services.CommentService, validate *validator.Validate) *CommentHandler {
	return &CommentHandler{svc: svc, validate: validate}
}

func (h *CommentHandler) create(w http.ResponseWriter, r *http.Request) {
	var req commentRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		httpx.ValidationError(w, err)
		return
	}
	m, err := h.svc.Create(r.Context(), domain.CreateCommentInput{
		Body: req.Body,
		Approved: req.Approved,
		PostID: req.PostID,
		AuthorID: req.AuthorID,
	})
	if err != nil {
		writeCommentError(w, err)
		return
	}
	httpx.JSON(w, http.StatusCreated, toCommentResponse(m))
}

func (h *CommentHandler) list(w http.ResponseWriter, r *http.Request) {
	limit := httpx.QueryInt(r, "limit", 20)
	offset := httpx.QueryInt(r, "offset", 0)
	items, total, err := h.svc.List(r.Context(), int32(limit), int32(offset))
	if err != nil {
		writeCommentError(w, err)
		return
	}
	out := make([]commentResponse, 0, len(items))
	for _, m := range items {
		out = append(out, toCommentResponse(m))
	}
	httpx.JSON(w, http.StatusOK, map[string]any{
		"data":   out,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *CommentHandler) get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	m, err := h.svc.Get(r.Context(), id)
	if err != nil {
		writeCommentError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toCommentResponse(m))
}

func (h *CommentHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req commentRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		httpx.ValidationError(w, err)
		return
	}
	m, err := h.svc.Update(r.Context(), id, domain.UpdateCommentInput{
		Body: req.Body,
		Approved: req.Approved,
		PostID: req.PostID,
		AuthorID: req.AuthorID,
	})
	if err != nil {
		writeCommentError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toCommentResponse(m))
}

func (h *CommentHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeCommentError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// RegisterRoutes mounts the comment CRUD endpoints (all require auth).
func (h *CommentHandler) RegisterRoutes(r chi.Router, tokens *auth.TokenService) {
	r.Group(func(r chi.Router) {
		r.Use(tokens.RequireAuth)
		r.Post("/blog/comments", h.create)
		r.Get("/blog/comments", h.list)
		r.Get("/blog/comments/{id}", h.get)
		r.Put("/blog/comments/{id}", h.update)
		r.Delete("/blog/comments/{id}", h.delete)
	})
}

func writeCommentError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrCommentNotFound):
		httpx.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrCommentConflict):
		httpx.Error(w, http.StatusConflict, err.Error())
	default:
		httpx.Error(w, http.StatusInternalServerError, "internal server error")
	}
}
