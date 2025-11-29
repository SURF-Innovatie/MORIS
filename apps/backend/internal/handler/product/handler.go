package product

import (
	"encoding/json"
	"net/http"

	"github.com/SURF-Innovatie/MORIS/internal/api/productdto"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	productsvc "github.com/SURF-Innovatie/MORIS/internal/product"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	svc productsvc.Service
}

func NewHandler(svc productsvc.Service) *Handler {
	return &Handler{svc: svc}
}

// Create godoc
// @Summary Create a product
// @Description Creates a new product
// @Tags products
// @Accept json
// @Produce json
// @Param product body productdto.Request true "Product data"
// @Success 200 {object} productdto.Response
// @Failure 400 {string} string "invalid body"
// @Failure 500 {string} string "internal server error"
// @Router /products [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req productdto.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
	}

	p, err := h.svc.Create(r.Context(), entities.Product{
		Name:     req.Name,
		Language: req.Language,
		Type:     entities.ProductType(req.Type),
		DOI:      req.DOI,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, productdto.Response{
		ID:       p.Id,
		Name:     p.Name,
		DOI:      p.DOI,
		Language: p.Language,
		Type:     p.Type,
	})
}

// GetAll godoc
// @Summary List products
// @Description Returns all products
// @Tags products
// @Produce json
// @Success 200 {array} productdto.Response
// @Failure 500 {string} string "internal server error"
// @Router /products [get]
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := h.svc.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, products)
}

// Get godoc
// @Summary Get a product
// @Description Get a single product by ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID (UUID)"
// @Success 200 {object} productdto.Response
// @Failure 400 {string} string "invalid id"
// @Failure 500 {string} string "internal server error"
// @Router /products/{id} [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	p, err := h.svc.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, productdto.Response{
		ID:       p.Id,
		Name:     p.Name,
		DOI:      p.DOI,
		Language: p.Language,
		Type:     p.Type,
	})
}

// Update godoc
// @Summary Update a product
// @Description Update an existing product
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID (UUID)"
// @Param product body productdto.Request true "Product data"
// @Success 200 {object} productdto.Response
// @Failure 400 {string} string "invalid id or body"
// @Failure 500 {string} string "internal server error"
// @Router /products/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
	}

	var req productdto.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
	}
	p, err := h.svc.Update(r.Context(), id, entities.Product{
		Name:     req.Name,
		Language: req.Language,
		Type:     entities.ProductType(req.Type),
		DOI:      req.DOI,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, productdto.Response{
		ID:       p.Id,
		Name:     p.Name,
		DOI:      p.DOI,
		Language: p.Language,
		Type:     p.Type,
	})
}

// Delete godoc
// @Summary Delete a product
// @Description Delete a product by ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID (UUID)"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "invalid id"
// @Failure 500 {string} string "internal server error"
// @Router /products/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
