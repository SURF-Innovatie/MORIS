package product

import (
	"net/http"

	"github.com/SURF-Innovatie/MORIS/internal/api/productdto"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/handler/middleware"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
	productsvc "github.com/SURF-Innovatie/MORIS/internal/product"
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
	if !httputil.ReadJSON(w, r, &req) {
		return
	}

	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
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

	_ = httputil.WriteJSON(w, http.StatusOK, productdto.Response{
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
// @Success 200 {array} entities.Product
// @Failure 500 {string} string "internal server error"
// @Router /products [get]
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := h.svc.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = httputil.WriteJSON(w, http.StatusOK, products)
}

// GetMe godoc
// @Summary List products for current user
// @Description Returns products associated with the current user
// @Tags products
// @Produce json
// @Success 200 {array} productdto.Response
// @Failure 500 {string} string "internal server error"
// @Router /products/me [get]
func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	userCtx, ok := middleware.GetUserFromContext(r.Context())
	if !ok || userCtx == nil {
		_ = httputil.WriteJSON(w, http.StatusUnauthorized, middleware.BackendError{
			Code:    http.StatusUnauthorized,
			Status:  "Unauthorized",
			Message: "User not authenticated or found in context",
		})
		return
	}

	products, err := h.svc.GetAllForUser(r.Context(), userCtx.User.PersonID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dtos := make([]productdto.Response, 0, len(products))
	for _, p := range products {
		dtos = append(dtos, productdto.Response{
			ID:       p.Id,
			Name:     p.Name,
			DOI:      p.DOI,
			Language: p.Language,
			Type:     p.Type,
		})
	}

	_ = httputil.WriteJSON(w, http.StatusOK, dtos)
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
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	p, err := h.svc.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = httputil.WriteJSON(w, http.StatusOK, productdto.Response{
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
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req productdto.Request
	if !httputil.ReadJSON(w, r, &req) {
		return
	}
	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
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
	_ = httputil.WriteJSON(w, http.StatusOK, productdto.Response{
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
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
