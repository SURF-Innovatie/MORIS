package project

import (
	"net/http"

	"github.com/SURF-Innovatie/MORIS/internal/api/dto"
	appauth "github.com/SURF-Innovatie/MORIS/internal/app/auth"
	"github.com/SURF-Innovatie/MORIS/internal/app/customfield"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/queries"
	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
	"github.com/SURF-Innovatie/MORIS/pkg/odrl"
	"github.com/samber/lo"
)

type Handler struct {
	svc            queries.Service
	customFieldSvc customfield.Service
	currentUser    appauth.CurrentUserProvider
}

func NewHandler(svc queries.Service, cfs customfield.Service, cu appauth.CurrentUserProvider) *Handler {
	return &Handler{svc: svc, customFieldSvc: cfs, currentUser: cu}
}

// GetProject godoc
// @Summary Get a project by ID
// @Description Retrieves a single project by its unique identifier
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID (UUID)"
// @Success 200 {object} dto.ProjectResponse
// @Failure 400 {string} string "invalid project id"
// @Failure 404 {string} string "project not found"
// @Router /projects/{id} [get]
func (h *Handler) GetProject(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid project id", nil)
		return
	}
	proj, err := h.svc.GetProject(r.Context(), id)
	if err != nil {
		httputil.WriteError(w, r, http.StatusNotFound, err.Error(), nil)
		return
	}
	_ = httputil.WriteJSON(w, http.StatusOK, transform.ToDTOItem[dto.ProjectResponse](proj))
}

// GetAllProjects godoc
// @Summary Get all projects
// @Description Retrieves a list of all projects
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.ProjectResponse
// @Failure 500 {string} string "internal server error"
// @Router /projects [get]
func (h *Handler) GetAllProjects(w http.ResponseWriter, r *http.Request) {

	projs, err := h.svc.GetAllProjects(r.Context())
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resps := transform.ToDTOs[dto.ProjectResponse](projs)
	_ = httputil.WriteJSON(w, http.StatusOK, resps)
}

// GetChangelog godoc
// @Summary Get change log for a project
// @Description Retrieves the change log for a specific project
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID (UUID)"
// @Success 200 {object} dto.Changelog
// @Failure 400 {string} string "invalid project id"
// @Failure 500 {string} string "internal server error"
// @Router /projects/{id}/changelog [get]
func (h *Handler) GetChangelog(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid project id", nil)
		return
	}

	changeLog, err := h.svc.GetChangeLog(r.Context(), id)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, transform.ToDTOItem[dto.Changelog](*changeLog))
}

// GetPendingEvents godoc
// @Summary Get pending events for a project
// @Description Retrieves a list of pending events for a specific project
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID (UUID)"
// @Success 200 {object} dto.EventResponse
// @Failure 400 {string} string "invalid project id"
// @Failure 500 {string} string "internal server error"
// @Router /projects/{id}/pending-events [get]
func (h *Handler) GetPendingEvents(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid project id", nil)
		return
	}

	pendingEvents, err := h.svc.GetPendingEvents(r.Context(), id)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp := dto.EventResponse{
		Events: lo.Map(pendingEvents, func(e events.DetailedEvent, _ int) dto.Event {
			var d dto.Event
			return d.FromDetailedEntity(e)
		}),
	}

	_ = httputil.WriteJSON(w, http.StatusOK, resp)
}

// ListAvailableRoles godoc
// @Summary List available roles for a project
// @Description Retrieves all roles available to be assigned in a project (inherited from organisation hierarchy)
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID (UUID)"
// @Success 200 {array} dto.ProjectRoleResponse
// @Failure 400 {string} string "invalid project id"
// @Failure 500 {string} string "internal server error"
// @Router /projects/{id}/roles [get]
func (h *Handler) ListAvailableRoles(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid project id", nil)
		return
	}

	roles, err := h.svc.ListAvailableRoles(r.Context(), id)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, transform.ToDTOs[dto.ProjectRoleResponse](roles))
}

// GetAllowedEvents godoc
// @Summary Get allowed events for a project
// @Description Retrieves a list of events the user is allowed to perform on the project
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID (UUID)"
// @Success 200 {array} string
// @Failure 400 {string} string "invalid project id"
// @Failure 500 {string} string "internal server error"
// @Router /projects/{id}/allowed-events [get]
func (h *Handler) GetAllowedEvents(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid project id", nil)
		return
	}

	allowedEvents, err := h.svc.GetAllowedEventTypes(r.Context(), id)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if len(allowedEvents) == 0 {
		// force write of empty array otherwise httputil.WriteJSON will write null
		httputil.WriteJSON(w, http.StatusOK, []string{})
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, allowedEvents)
}

// ListAvailableCustomFields godoc
// @Summary List available custom fields for a project
// @Description Retrieves all custom fields available to be populated in a project (inherited from organisation hierarchy)
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID (UUID)"
// @Success 200 {array} dto.CustomFieldDefinitionResponse
// @Failure 400 {string} string "invalid project id"
// @Failure 500 {string} string "internal server error"
// @Router /projects/{id}/custom-fields [get]
func (h *Handler) ListAvailableCustomFields(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid project id", nil)
		return
	}

	// 1. Get Project to find OwningOrgNode
	proj, err := h.svc.GetProject(r.Context(), id)
	if err != nil {
		httputil.WriteError(w, r, http.StatusNotFound, "project not found", nil)
		return
	}

	// 2. List definitions for that Org Node
	category := entities.CustomFieldCategory("PROJECT")
	defs, err := h.customFieldSvc.ListAvailableForNode(r.Context(), proj.OwningOrgNode.ID, &category)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, transform.ToDTOs[dto.CustomFieldDefinitionResponse](defs))
}

// GetODRLPolicy godoc
// @Summary Get ODRL policy for a project
// @Description Retrieves an ODRL policy document expressing the current user's permissions in the project
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID (UUID)"
// @Success 200 {object} odrl.Policy
// @Failure 400 {string} string "invalid project id"
// @Failure 500 {string} string "internal server error"
// @Router /projects/{id}/odrl-policy [get]
func (h *Handler) GetODRLPolicy(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid project id", nil)
		return
	}

	// Get user's allowed event types for this project
	allowedEvents, err := h.svc.GetAllowedEventTypes(r.Context(), id)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Get current user ID for the assignee
	u, err := h.currentUser.Current(r.Context())
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	userID := u.UserID()

	// Build the ODRL policy
	projectURI := "urn:moris:project:" + id.String()
	userURI := "urn:moris:user:" + userID.String()

	policy := odrl.NewSet(projectURI + ":policy")
	for _, eventType := range allowedEvents {
		policy.AddPermission(odrl.NewPermission(eventType).
			WithTarget(projectURI).
			WithAssignee(userURI))
	}

	_ = httputil.WriteJSON(w, http.StatusOK, policy)
}
