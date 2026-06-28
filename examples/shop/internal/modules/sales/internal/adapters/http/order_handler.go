// Package saleshttp is the sales module's HTTP transport layer.
package saleshttp

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/example/shop/internal/modules/sales/internal/domain"
	"github.com/example/shop/internal/modules/sales/internal/services"
	"github.com/example/shop/internal/platform/auth"
	"github.com/example/shop/internal/platform/httpx"
)

type orderRequest struct {
	Status string `json:"status" validate:"required"`
	Total string `json:"total" validate:"required"`
	PlacedAt time.Time `json:"placed_at"`
	CustomerID uuid.UUID `json:"customer_id" validate:"required"`
}

type orderResponse struct {
	ID string `json:"id"`
	Status string `json:"status"`
	Total string `json:"total"`
	PlacedAt time.Time `json:"placed_at"`
	CustomerID uuid.UUID `json:"customer_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toOrderResponse(m domain.Order) orderResponse {
	return orderResponse{
		ID: m.ID.String(),
		Status: m.Status,
		Total: m.Total,
		PlacedAt: m.PlacedAt,
		CustomerID: m.CustomerID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// OrderHandler exposes CRUD endpoints for orders.
type OrderHandler struct {
	svc      *services.OrderService
	validate *validator.Validate
}

// NewOrderHandler constructs a OrderHandler.
func NewOrderHandler(svc *services.OrderService, validate *validator.Validate) *OrderHandler {
	return &OrderHandler{svc: svc, validate: validate}
}

func (h *OrderHandler) create(w http.ResponseWriter, r *http.Request) {
	var req orderRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		httpx.ValidationError(w, err)
		return
	}
	m, err := h.svc.Create(r.Context(), domain.CreateOrderInput{
		Status: req.Status,
		Total: req.Total,
		PlacedAt: req.PlacedAt,
		CustomerID: req.CustomerID,
	})
	if err != nil {
		writeOrderError(w, err)
		return
	}
	httpx.JSON(w, http.StatusCreated, toOrderResponse(m))
}

func (h *OrderHandler) list(w http.ResponseWriter, r *http.Request) {
	limit := httpx.QueryInt(r, "limit", 20)
	offset := httpx.QueryInt(r, "offset", 0)
	items, total, err := h.svc.List(r.Context(), int32(limit), int32(offset))
	if err != nil {
		writeOrderError(w, err)
		return
	}
	out := make([]orderResponse, 0, len(items))
	for _, m := range items {
		out = append(out, toOrderResponse(m))
	}
	httpx.JSON(w, http.StatusOK, map[string]any{
		"data":   out,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *OrderHandler) get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	m, err := h.svc.Get(r.Context(), id)
	if err != nil {
		writeOrderError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toOrderResponse(m))
}

func (h *OrderHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req orderRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		httpx.ValidationError(w, err)
		return
	}
	m, err := h.svc.Update(r.Context(), id, domain.UpdateOrderInput{
		Status: req.Status,
		Total: req.Total,
		PlacedAt: req.PlacedAt,
		CustomerID: req.CustomerID,
	})
	if err != nil {
		writeOrderError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toOrderResponse(m))
}

func (h *OrderHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeOrderError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// RegisterRoutes mounts the order CRUD endpoints (all require auth).
func (h *OrderHandler) RegisterRoutes(r chi.Router, tokens *auth.TokenService) {
	r.Group(func(r chi.Router) {
		r.Use(tokens.RequireAuth)
		r.Post("/sales/orders", h.create)
		r.Get("/sales/orders", h.list)
		r.Get("/sales/orders/{id}", h.get)
		r.Put("/sales/orders/{id}", h.update)
		r.Delete("/sales/orders/{id}", h.delete)
	})
}

func writeOrderError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrOrderNotFound):
		httpx.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrOrderConflict):
		httpx.Error(w, http.StatusConflict, err.Error())
	default:
		httpx.Error(w, http.StatusInternalServerError, "internal server error")
	}
}
