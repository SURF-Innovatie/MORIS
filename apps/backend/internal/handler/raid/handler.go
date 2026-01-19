package raid

import (
	"fmt"
	"net/http"

	"github.com/SURF-Innovatie/MORIS/external/raid"
	"github.com/SURF-Innovatie/MORIS/internal/api/dto"
	appauth "github.com/SURF-Innovatie/MORIS/internal/app/auth"
	appraid "github.com/SURF-Innovatie/MORIS/internal/app/raid"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
	"github.com/go-chi/chi/v5"
)

// Handler provides HTTP handlers for RAiD integration
type Handler struct {
	raidService   appraid.Service
	mapperService appraid.MapperService
	currentUser   appauth.CurrentUserProvider
}

// NewHandler creates a new RAiD handler
func NewHandler(raidService appraid.Service, mapperService appraid.MapperService, currentUser appauth.CurrentUserProvider) *Handler {
	return &Handler{
		raidService:   raidService,
		mapperService: mapperService,
		currentUser:   currentUser,
	}
}

// MountRoutes mounts the RAiD routes on the router
func MountRoutes(r chi.Router, h *Handler) {
	// r is expected to be mounted at /projects/{id}
	r.Get("/raid", h.GetProjectRaid)
	r.Post("/raid", h.MintProjectRaid)
	r.Put("/raid", h.UpdateProjectRaid)
}

// GetProjectRaid godoc
// @Summary Get RAiD for a project
// @Description Returns the RAiD associated with the project
// @Tags raid
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Success 200 {object} dto.RAiD
// @Failure 404 {object} httputil.BackendError "RAiD not found"
// @Failure 401 {object} httputil.BackendError "User not authenticated"
// @Failure 500 {object} httputil.BackendError "Internal server error"
// @Router /projects/{id}/raid [get]
func (h *Handler) GetProjectRaid(w http.ResponseWriter, r *http.Request) {
	// Check auth
	_, err := h.currentUser.Current(r.Context())
	if err != nil {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	projectID, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid project id", nil)
		return
	}

	info, err := h.raidService.GetRaidInfoByProject(r.Context(), projectID)
	if err != nil {
		// If not found locally, we return 404
		httputil.WriteError(w, r, http.StatusNotFound, "raid not found", nil)
		return
	}
	_ = info

	httputil.WriteError(w, r, http.StatusNotImplemented, "get project raid not fully implemented", nil)
}

// MintProjectRaid godoc
// @Summary Mint a new RAiD for a project
// @Description Mints a new RAiD based on project data
// @Tags raid
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Success 201 {object} dto.RAiD
// @Failure 400 {object} httputil.BackendError "Invalid request or compatibility issues"
// @Failure 401 {object} httputil.BackendError "User not authenticated"
// @Failure 500 {object} httputil.BackendError "Internal server error"
// @Router /projects/{id}/raid [post]
func (h *Handler) MintProjectRaid(w http.ResponseWriter, r *http.Request) {
	u, err := h.currentUser.Current(r.Context())
	if err != nil {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	projectID, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid project id", nil)
		return
	}

	createReq, err := h.mapperService.MapToCreateRequest(r.Context(), projectID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// MintRaid now handles linking to project and saving local info
	raidDto, err := h.raidService.MintRaid(r.Context(), u.UserID(), projectID, createReq)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	httputil.WriteJSON(w, http.StatusCreated, dto.FromExternalRaid(*raidDto))
}

// UpdateProjectRaid godoc
// @Summary Update an existing RAiD for a project
// @Description Updates the RAiD associated with the project
// @Tags raid
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Success 200 {object} dto.RAiD
// @Failure 400 {object} httputil.BackendError "Invalid request"
// @Failure 404 {object} httputil.BackendError "RAiD not found"
// @Failure 401 {object} httputil.BackendError "User not authenticated"
// @Failure 500 {object} httputil.BackendError "Internal server error"
// @Router /projects/{id}/raid [put]
func (h *Handler) UpdateProjectRaid(w http.ResponseWriter, r *http.Request) {
	u, err := h.currentUser.Current(r.Context())
	if err != nil {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	projectID, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid project id", nil)
		return
	}

	// Find local RAiD info to get identifier
	info, err := h.raidService.GetRaidInfoByProject(r.Context(), projectID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusNotFound, "raid not found for project", nil)
		return
	}

	// Construct RAiDId from local info
	var servicePoint *string
	if info.OwnerServicePoint != nil {
		sp := fmt.Sprintf("%d", *info.OwnerServicePoint)
		servicePoint = &sp
	}

	identifier := raid.RAiDId{
		IdValue:   info.RAiDId,
		SchemaUri: info.SchemaUri,
		RegistrationAgency: raid.RAiDRegistrationAgency{
			Id:        info.RegistrationAgencyId,
			SchemaUri: info.RegistrationAgencySchemaUri,
		},
		Owner: raid.RAiDOwner{
			Id:           info.OwnerId,
			SchemaUri:    info.OwnerSchemaUri,
			ServicePoint: servicePoint,
		},
		License: info.License,
		Version: info.Version,
	}

	// Use MapToUpdateRequest to prepare the request
	updateReq, err := h.mapperService.MapToUpdateRequest(r.Context(), projectID, identifier)
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// UpdateRaid now handles parsing handle and saving local info
	raidDto, err := h.raidService.UpdateRaid(r.Context(), u.UserID(), projectID, info.RAiDId, updateReq)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, dto.FromExternalRaid(*raidDto))
}
