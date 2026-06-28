// Package cataloghttp is the catalog module's HTTP transport layer.
package cataloghttp

import (
	"encoding/json"
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

type productRequest struct {
	Name string `json:"name" validate:"required"`
	Sku string `json:"sku" validate:"required"`
	Price string `json:"price" validate:"required"`
	Stock int32 `json:"stock"`
	Active bool `json:"active"`
	Metadata json.RawMessage `json:"metadata"`
	CategoryID uuid.UUID `json:"category_id" validate:"required"`
}

type productResponse struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Sku string `json:"sku"`
	Price string `json:"price"`
	Stock int32 `json:"stock"`
	Active bool `json:"active"`
	Metadata json.RawMessage `json:"metadata"`
	CategoryID uuid.UUID `json:"category_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toProductResponse(m domain.Product) productResponse {
	return productResponse{
		ID: m.ID.String(),
		Name: m.Name,
		Sku: m.Sku,
		Price: m.Price,
		Stock: m.Stock,
		Active: m.Active,
		Metadata: m.Metadata,
		CategoryID: m.CategoryID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// ProductHandler exposes CRUD endpoints for products.
type ProductHandler struct {
	svc      *services.ProductService
	validate *validator.Validate
}

// NewProductHandler constructs a ProductHandler.
func NewProductHandler(svc *services.ProductService, validate *validator.Validate) *ProductHandler {
	return &ProductHandler{svc: svc, validate: validate}
}

func (h *ProductHandler) create(w http.ResponseWriter, r *http.Request) {
	var req productRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		httpx.ValidationError(w, err)
		return
	}
	m, err := h.svc.Create(r.Context(), domain.CreateProductInput{
		Name: req.Name,
		Sku: req.Sku,
		Price: req.Price,
		Stock: req.Stock,
		Active: req.Active,
		Metadata: req.Metadata,
		CategoryID: req.CategoryID,
	})
	if err != nil {
		writeProductError(w, err)
		return
	}
	httpx.JSON(w, http.StatusCreated, toProductResponse(m))
}

func (h *ProductHandler) list(w http.ResponseWriter, r *http.Request) {
	limit := httpx.QueryInt(r, "limit", 20)
	offset := httpx.QueryInt(r, "offset", 0)
	items, total, err := h.svc.List(r.Context(), int32(limit), int32(offset))
	if err != nil {
		writeProductError(w, err)
		return
	}
	out := make([]productResponse, 0, len(items))
	for _, m := range items {
		out = append(out, toProductResponse(m))
	}
	httpx.JSON(w, http.StatusOK, map[string]any{
		"data":   out,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *ProductHandler) get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	m, err := h.svc.Get(r.Context(), id)
	if err != nil {
		writeProductError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toProductResponse(m))
}

func (h *ProductHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req productRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		httpx.ValidationError(w, err)
		return
	}
	m, err := h.svc.Update(r.Context(), id, domain.UpdateProductInput{
		Name: req.Name,
		Sku: req.Sku,
		Price: req.Price,
		Stock: req.Stock,
		Active: req.Active,
		Metadata: req.Metadata,
		CategoryID: req.CategoryID,
	})
	if err != nil {
		writeProductError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toProductResponse(m))
}

func (h *ProductHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeProductError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// RegisterRoutes mounts the product CRUD endpoints (all require auth).
func (h *ProductHandler) RegisterRoutes(r chi.Router, tokens *auth.TokenService) {
	r.Group(func(r chi.Router) {
		r.Use(tokens.RequireAuth)
		r.Post("/catalog/products", h.create)
		r.Get("/catalog/products", h.list)
		r.Get("/catalog/products/{id}", h.get)
		r.Put("/catalog/products/{id}", h.update)
		r.Delete("/catalog/products/{id}", h.delete)
	})
}

func writeProductError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrProductNotFound):
		httpx.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrProductConflict):
		httpx.Error(w, http.StatusConflict, err.Error())
	default:
		httpx.Error(w, http.StatusInternalServerError, "internal server error")
	}
}
