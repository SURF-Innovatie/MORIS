package raid

import (
	"net/http"

	"github.com/SURF-Innovatie/MORIS/internal/api/dto"
	appauth "github.com/SURF-Innovatie/MORIS/internal/app/auth"
	appraid "github.com/SURF-Innovatie/MORIS/internal/app/raid"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
	"github.com/go-chi/chi/v5"
)

// Handler provides HTTP handlers for RAiD integration
type Handler struct {
	raidService appraid.Service
	currentUser appauth.CurrentUserProvider
}

// NewHandler creates a new RAiD handler
func NewHandler(raidService appraid.Service, currentUser appauth.CurrentUserProvider) *Handler {
	return &Handler{
		raidService: raidService,
		currentUser: currentUser,
	}
}

// MountRoutes mounts the RAiD routes on the router
func MountRoutes(r chi.Router, h *Handler) {
	r.Get("/raid", h.ListRaids)
	r.Post("/raid", h.MintRaid)
	r.Get("/raid/{prefix}/{suffix}", h.GetRaid)
	r.Put("/raid/{prefix}/{suffix}", h.UpdateRaid)
}

// ListRaids godoc
// @Summary List all RAiDs
// @Description Returns all RAiDs
// @Tags raid
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.RAiD
// @Failure 401 {object} httputil.BackendError "User not authenticated"
// @Failure 500 {object} httputil.BackendError "Internal server error"
// @Router /raid [get]
func (h *Handler) ListRaids(w http.ResponseWriter, r *http.Request) {
	u, err := h.currentUser.Current(r.Context())
	if err != nil {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	raids, err := h.raidService.FindAllRaids(r.Context(), u.UserID())
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, dto.FromExternalRaids(raids))
}

// MintRaid godoc
// @Summary Mint a new RAiD
// @Description Mints a new RAiD
// @Tags raid
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.RAiDCreateRequest true "RAiD create request"
// @Success 201 {object} dto.RAiD
// @Failure 400 {object} httputil.BackendError "Invalid request"
// @Failure 401 {object} httputil.BackendError "User not authenticated"
// @Failure 500 {object} httputil.BackendError "Internal server error"
// @Router /raid [post]
func (h *Handler) MintRaid(w http.ResponseWriter, r *http.Request) {
	u, err := h.currentUser.Current(r.Context())
	if err != nil {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	var req dto.RAiDCreateRequest
	if !httputil.ReadJSON(w, r, &req) {
		return
	}

	raid, err := h.raidService.MintRaid(r.Context(), u.UserID(), &req)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	httputil.WriteJSON(w, http.StatusCreated, dto.FromExternalRaid(*raid))
}

// GetRaid godoc
// @Summary Get a RAiD
// @Description Returns a single RAiD by handle
// @Tags raid
// @Produce json
// @Security BearerAuth
// @Param prefix path string true "RAiD Prefix"
// @Param suffix path string true "RAiD Suffix"
// @Success 200 {object} dto.RAiD
// @Failure 400 {object} httputil.BackendError "Invalid handle"
// @Failure 401 {object} httputil.BackendError "User not authenticated"
// @Failure 500 {object} httputil.BackendError "Internal server error"
// @Router /raid/{prefix}/{suffix} [get]
func (h *Handler) GetRaid(w http.ResponseWriter, r *http.Request) {
	u, err := h.currentUser.Current(r.Context())
	if err != nil {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	prefix := chi.URLParam(r, "prefix")
	suffix := chi.URLParam(r, "suffix")
	if prefix == "" || suffix == "" {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid raid handle", nil)
		return
	}

	raid, err := h.raidService.FindRaid(r.Context(), u.UserID(), prefix, suffix)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, dto.FromExternalRaid(*raid))
}

// UpdateRaid godoc
// @Summary Update a RAiD
// @Description Updates an existing RAiD
// @Tags raid
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param prefix path string true "RAiD Prefix"
// @Param suffix path string true "RAiD Suffix"
// @Param body body dto.RAiDUpdateRequest true "RAiD update request"
// @Success 200 {object} dto.RAiD
// @Failure 400 {object} httputil.BackendError "Invalid request"
// @Failure 401 {object} httputil.BackendError "User not authenticated"
// @Failure 500 {object} httputil.BackendError "Internal server error"
// @Router /raid/{prefix}/{suffix} [put]
func (h *Handler) UpdateRaid(w http.ResponseWriter, r *http.Request) {
	u, err := h.currentUser.Current(r.Context())
	if err != nil {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	prefix := chi.URLParam(r, "prefix")
	suffix := chi.URLParam(r, "suffix")
	if prefix == "" || suffix == "" {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid raid handle", nil)
		return
	}

	var req dto.RAiDUpdateRequest
	if !httputil.ReadJSON(w, r, &req) {
		return
	}

	raid, err := h.raidService.UpdateRaid(r.Context(), u.UserID(), prefix, suffix, &req)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, dto.FromExternalRaid(*raid))
}
