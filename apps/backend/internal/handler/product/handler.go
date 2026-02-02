package product

import (
	"net/http"

	"github.com/SURF-Innovatie/MORIS/internal/api/dto"
	appauth "github.com/SURF-Innovatie/MORIS/internal/app/auth"
	productsvc "github.com/SURF-Innovatie/MORIS/internal/app/product"
	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
)

type Handler struct {
	svc         productsvc.Service
	currentUser appauth.CurrentUserProvider
}

func NewHandler(svc productsvc.Service, cu appauth.CurrentUserProvider) *Handler {
	return &Handler{svc: svc, currentUser: cu}
}

// Create godoc
// @Summary Create a product
// @Description Creates a new product
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param product body dto.ProductRequest true "Product data"
// @Success 200 {object} dto.ProductResponse
// @Failure 400 {string} string "invalid body"
// @Failure 500 {string} string "internal server error"
// @Router /products [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.ProductRequest
	if !httputil.ReadJSON(w, r, &req) {
		return
	}

	if req.Name == "" {
		httputil.WriteError(w, r, http.StatusBadRequest, "name is required", nil)
		return
	}

	// Get the current user to attribute the product to them
	u, err := h.currentUser.Current(r.Context())
	if err != nil {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	var zenodoID int
	if req.ZenodoDepositionID != nil {
		zenodoID = *req.ZenodoDepositionID
	}

	p, err := h.svc.Create(r.Context(), entities.Product{
		Name:               req.Name,
		Language:           req.Language,
		Type:               entities.ProductType(req.Type),
		DOI:                req.DOI,
		ZenodoDepositionID: zenodoID,
		AuthorPersonID:     u.PersonID,
	})
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, transform.ToDTOItem[dto.ProductResponse](*p))
}

// GetAll godoc
// @Summary List products
// @Description Returns all products
// @Tags products
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.ProductResponse
// @Failure 500 {string} string "internal server error"
// @Router /products [get]
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := h.svc.GetAll(r.Context())
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	dtos := make([]dto.ProductResponse, len(products))
	for i, p := range products {
		dtos[i] = transform.ToDTOItem[dto.ProductResponse](*p)
	}

	_ = httputil.WriteJSON(w, http.StatusOK, dtos)
}

// GetMe godoc
// @Summary List products for current user
// @Description Returns products associated with the current user
// @Tags products
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.ProductResponse
// @Failure 500 {string} string "internal server error"
// @Router /products/me [get]
func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	u, err := h.currentUser.Current(r.Context())
	if err != nil {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	products, err := h.svc.GetAllForUser(r.Context(), u.PersonID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// TODO: make getall return non pointer slice
	dtos := make([]dto.ProductResponse, len(products))
	for i, p := range products {
		dtos[i] = transform.ToDTOItem[dto.ProductResponse](*p)
	}

	_ = httputil.WriteJSON(w, http.StatusOK, dtos)
}

// Get godoc
// @Summary Get a product
// @Description Get a single product by ID
// @Tags products
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID (UUID)"
// @Success 200 {object} dto.ProductResponse
// @Failure 400 {string} string "invalid id"
// @Failure 500 {string} string "internal server error"
// @Router /products/{id} [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid id", nil)
		return
	}
	p, err := h.svc.Get(r.Context(), id)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	_ = httputil.WriteJSON(w, http.StatusOK, transform.ToDTOItem[dto.ProductResponse](*p))
}

// Update godoc
// @Summary Update a product
// @Description Update an existing product
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID (UUID)"
// @Param product body dto.ProductRequest true "Product data"
// @Success 200 {object} dto.ProductResponse
// @Failure 400 {string} string "invalid id or body"
// @Failure 500 {string} string "internal server error"
// @Router /products/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid id", nil)
		return
	}

	var req dto.ProductRequest
	if !httputil.ReadJSON(w, r, &req) {
		return
	}
	if req.Name == "" {
		httputil.WriteError(w, r, http.StatusBadRequest, "name is required", nil)
		return
	}
	var zenodoID int
	if req.ZenodoDepositionID != nil {
		zenodoID = *req.ZenodoDepositionID
	}
	p, err := h.svc.Update(r.Context(), id, entities.Product{
		Name:               req.Name,
		Language:           req.Language,
		Type:               entities.ProductType(req.Type),
		DOI:                req.DOI,
		ZenodoDepositionID: zenodoID,
	})
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	_ = httputil.WriteJSON(w, http.StatusOK, transform.ToDTOItem[dto.ProductResponse](*p))
}

// Delete godoc
// @Summary Delete a product
// @Description Delete a product by ID
// @Tags products
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID (UUID)"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "invalid id"
// @Failure 500 {string} string "internal server error"
// @Router /products/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid id", nil)
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, "failed to delete product", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
