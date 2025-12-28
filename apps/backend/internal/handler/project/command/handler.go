package command

import (
	"net/http"

	"github.com/SURF-Innovatie/MORIS/internal/api/dto"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
	"github.com/SURF-Innovatie/MORIS/internal/project/command"
)

type Handler struct {
	svc command.Service
}

func NewHandler(svc command.Service) *Handler {
	return &Handler{svc: svc}
}

// ListAvailableEvents godoc
// @Summary List available events for a project
// @Description Lists all event types that can be executed for a project
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID (UUID)"
// @Success 200 {array} dto.AvailableEvent
// @Failure 400 {string} string "invalid project id"
// @Failure 500 {string} string "internal server error"
// @Router /projects/{id}/events [get]
func (h *Handler) ListAvailableEvents(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid project id", nil)
		return
	}

	events, err := h.svc.ListAvailableEvents(r.Context(), &id)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, events)
}

// ExecuteEvent godoc
// @Summary Execute a project event
// @Description Executes a single event against a project
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID (UUID)"
// @Param body body dto.ExecuteEventRequest true "Event execution request"
// @Success 200 {object} object "Updated project"
// @Failure 400 {string} string "invalid request"
// @Failure 404 {string} string "unknown event type"
// @Failure 500 {string} string "internal server error"
// @Router /projects/{id}/events [post]
func (h *Handler) ExecuteEvent(w http.ResponseWriter, r *http.Request) {
	projectID, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid project id", nil)
		return
	}

	var req dto.ExecuteEventRequest
	if !httputil.ReadJSON(w, r, &req) {
		return
	}

	// enforce path ID
	req.ProjectID = projectID

	proj, err := h.svc.ExecuteEvent(r.Context(), req)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, proj)
}
