package catalog

import (
	"errors"
	"net/http"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Create creates a new catalog.
// @Summary Create a new catalog
// @Description Create a new catalog with the given name, title, and optional fields.
// @Tags catalogs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateRequest true "Create Catalog Request"
// @Success 201 {object} ent.Catalog
// @Failure 400 {object} httputil.BackendError
// @Failure 401 {object} httputil.BackendError
// @Failure 403 {object} httputil.BackendError
// @Failure 500 {object} httputil.BackendError
// @Router /catalogs [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if !httputil.ReadJSON(w, r, &req) {
		return
	}

	if req.Name == "" || req.Title == "" {
		httputil.WriteError(w, r, http.StatusBadRequest, "name and title are required", nil)
		return
	}

	cat, err := h.service.Create(r.Context(), req)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, "failed to create catalog", err.Error())
		return
	}

	httputil.WriteJSON(w, http.StatusCreated, cat)
}

// Get retrieves a catalog by ID with all aggregated details.
// @Summary Get a catalog by ID
// @Description Get a catalog by ID including all aggregated project details.
// @Tags catalogs
// @Produce json
// @Param id path string true "Catalog ID"
// @Success 200 {object} CatalogDetails
// @Failure 400 {object} httputil.BackendError
// @Failure 404 {object} httputil.BackendError
// @Failure 500 {object} httputil.BackendError
// @Router /catalogs/{id} [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid catalog id", err.Error())
		return
	}

	details, err := h.service.GetDetails(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) || ent.IsNotFound(err) {
			httputil.WriteError(w, r, http.StatusNotFound, "catalog not found", nil)
			return
		}
		httputil.WriteError(w, r, http.StatusInternalServerError, "failed to get catalog", err.Error())
		return
	}

	httputil.WriteJSON(w, http.StatusOK, details)
}

// List lists all catalogs.
// @Summary List all catalogs
// @Description List all catalogs.
// @Tags catalogs
// @Produce json
// @Success 200 {array} ent.Catalog
// @Failure 500 {object} httputil.BackendError
// @Router /catalogs [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	cats, err := h.service.List(r.Context())
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, "failed to list catalogs", err.Error())
		return
	}

	httputil.WriteJSON(w, http.StatusOK, cats)
}

// Update updates an existing catalog.
// @Summary Update a catalog
// @Description Update a catalog by ID. Only provided fields will be updated.
// @Tags catalogs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Catalog ID"
// @Param request body UpdateRequest true "Update Catalog Request"
// @Success 200 {object} ent.Catalog
// @Failure 400 {object} httputil.BackendError
// @Failure 401 {object} httputil.BackendError
// @Failure 403 {object} httputil.BackendError
// @Failure 404 {object} httputil.BackendError
// @Failure 500 {object} httputil.BackendError
// @Router /catalogs/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid catalog id", err.Error())
		return
	}

	var req UpdateRequest
	if !httputil.ReadJSON(w, r, &req) {
		return
	}

	cat, err := h.service.Update(r.Context(), id, req)
	if err != nil {
		if ent.IsNotFound(err) {
			httputil.WriteError(w, r, http.StatusNotFound, "catalog not found", nil)
			return
		}
		httputil.WriteError(w, r, http.StatusInternalServerError, "failed to update catalog", err.Error())
		return
	}

	httputil.WriteJSON(w, http.StatusOK, cat)
}

// Delete deletes a catalog by ID.
// @Summary Delete a catalog
// @Description Delete a catalog by ID.
// @Tags catalogs
// @Produce json
// @Security BearerAuth
// @Param id path string true "Catalog ID"
// @Success 204
// @Failure 400 {object} httputil.BackendError
// @Failure 401 {object} httputil.BackendError
// @Failure 403 {object} httputil.BackendError
// @Failure 404 {object} httputil.BackendError
// @Failure 500 {object} httputil.BackendError
// @Router /catalogs/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid catalog id", err.Error())
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if ent.IsNotFound(err) {
			httputil.WriteError(w, r, http.StatusNotFound, "catalog not found", nil)
			return
		}
		httputil.WriteError(w, r, http.StatusInternalServerError, "failed to delete catalog", err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// MountRoutes mounts catalog routes on the given router.
// GET endpoints are public; POST, PUT, DELETE require auth + sysadmin middleware.
func MountRoutes(r chi.Router, h *Handler, authMiddleware, adminMiddleware func(http.Handler) http.Handler) {
	r.Route("/catalogs", func(r chi.Router) {
		// Public
		r.Get("/", h.List)
		r.Get("/{id}", h.Get)

		// Admin only
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)
			r.Use(adminMiddleware)
			r.Post("/", h.Create)
			r.Put("/{id}", h.Update)
			r.Delete("/{id}", h.Delete)
		})
	})
}
