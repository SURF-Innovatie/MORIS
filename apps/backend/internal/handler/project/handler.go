package project

import (
	"net/http"

	"github.com/SURF-Innovatie/MORIS/internal/api/dto"
	"github.com/SURF-Innovatie/MORIS/internal/common/transform"

	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
	"github.com/SURF-Innovatie/MORIS/internal/project"
)

type Handler struct {
	svc project.Service
}

func NewHandler(svc project.Service) *Handler {
	return &Handler{svc: svc}
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
		Events: make([]dto.Event, 0, len(pendingEvents)),
	}
	for _, e := range pendingEvents {
		resp.Events = append(resp.Events, transform.ToDTOItem[dto.Event](e))
	}

	_ = httputil.WriteJSON(w, http.StatusOK, resp)
}

// GetProjectRoles godoc
// @Summary Get project roles
// @Description Retrieves a list of all available project roles
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.ProjectRoleResponse
// @Failure 500 {string} string "internal server error"
// @Router /projects/roles [get]
func (h *Handler) GetProjectRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := h.svc.GetProjectRoles(r.Context())
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resps := make([]dto.ProjectRoleResponse, 0, len(roles))
	for _, role := range roles {
		resps = append(resps, dto.ProjectRoleResponse{
			ID:   role.ID,
			Key:  role.Key,
			Name: role.Name,
		})
	}

	_ = httputil.WriteJSON(w, http.StatusOK, resps)
}
