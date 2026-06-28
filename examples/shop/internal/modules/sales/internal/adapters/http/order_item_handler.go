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

type orderItemRequest struct {
	Quantity int32 `json:"quantity" validate:"required"`
	UnitPrice string `json:"unit_price" validate:"required"`
	OrderID uuid.UUID `json:"order_id" validate:"required"`
	ProductID uuid.UUID `json:"product_id" validate:"required"`
}

type orderItemResponse struct {
	ID string `json:"id"`
	Quantity int32 `json:"quantity"`
	UnitPrice string `json:"unit_price"`
	OrderID uuid.UUID `json:"order_id"`
	ProductID uuid.UUID `json:"product_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toOrderItemResponse(m domain.OrderItem) orderItemResponse {
	return orderItemResponse{
		ID: m.ID.String(),
		Quantity: m.Quantity,
		UnitPrice: m.UnitPrice,
		OrderID: m.OrderID,
		ProductID: m.ProductID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// OrderItemHandler exposes CRUD endpoints for order_items.
type OrderItemHandler struct {
	svc      *services.OrderItemService
	validate *validator.Validate
}

// NewOrderItemHandler constructs a OrderItemHandler.
func NewOrderItemHandler(svc *services.OrderItemService, validate *validator.Validate) *OrderItemHandler {
	return &OrderItemHandler{svc: svc, validate: validate}
}

func (h *OrderItemHandler) create(w http.ResponseWriter, r *http.Request) {
	var req orderItemRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		httpx.ValidationError(w, err)
		return
	}
	m, err := h.svc.Create(r.Context(), domain.CreateOrderItemInput{
		Quantity: req.Quantity,
		UnitPrice: req.UnitPrice,
		OrderID: req.OrderID,
		ProductID: req.ProductID,
	})
	if err != nil {
		writeOrderItemError(w, err)
		return
	}
	httpx.JSON(w, http.StatusCreated, toOrderItemResponse(m))
}

func (h *OrderItemHandler) list(w http.ResponseWriter, r *http.Request) {
	limit := httpx.QueryInt(r, "limit", 20)
	offset := httpx.QueryInt(r, "offset", 0)
	items, total, err := h.svc.List(r.Context(), int32(limit), int32(offset))
	if err != nil {
		writeOrderItemError(w, err)
		return
	}
	out := make([]orderItemResponse, 0, len(items))
	for _, m := range items {
		out = append(out, toOrderItemResponse(m))
	}
	httpx.JSON(w, http.StatusOK, map[string]any{
		"data":   out,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *OrderItemHandler) get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	m, err := h.svc.Get(r.Context(), id)
	if err != nil {
		writeOrderItemError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toOrderItemResponse(m))
}

func (h *OrderItemHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req orderItemRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		httpx.ValidationError(w, err)
		return
	}
	m, err := h.svc.Update(r.Context(), id, domain.UpdateOrderItemInput{
		Quantity: req.Quantity,
		UnitPrice: req.UnitPrice,
		OrderID: req.OrderID,
		ProductID: req.ProductID,
	})
	if err != nil {
		writeOrderItemError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toOrderItemResponse(m))
}

func (h *OrderItemHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeOrderItemError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// RegisterRoutes mounts the order_item CRUD endpoints (all require auth).
func (h *OrderItemHandler) RegisterRoutes(r chi.Router, tokens *auth.TokenService) {
	r.Group(func(r chi.Router) {
		r.Use(tokens.RequireAuth)
		r.Post("/sales/order-items", h.create)
		r.Get("/sales/order-items", h.list)
		r.Get("/sales/order-items/{id}", h.get)
		r.Put("/sales/order-items/{id}", h.update)
		r.Delete("/sales/order-items/{id}", h.delete)
	})
}

func writeOrderItemError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrOrderItemNotFound):
		httpx.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrOrderItemConflict):
		httpx.Error(w, http.StatusConflict, err.Error())
	default:
		httpx.Error(w, http.StatusInternalServerError, "internal server error")
	}
}
