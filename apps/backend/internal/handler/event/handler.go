package event

import (
	"net/http"

	"github.com/SURF-Innovatie/MORIS/internal/api/dto"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/SURF-Innovatie/MORIS/internal/event"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
)

type Handler struct {
	svc event.Service
}

func NewHandler(svc event.Service) *Handler {
	return &Handler{svc: svc}
}

// GetEvent godoc
// @Summary Get event details
// @Description Retrieves details for a specific event by ID
// @Tags events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Event ID (UUID)"
// @Success 200 {object} dto.Event
// @Failure 400 {string} string "invalid event id"
// @Failure 404 {string} string "event not found"
// @Failure 500 {string} string "internal server error"
// @Router /events/{id} [get]
func (h *Handler) GetEvent(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid event id", nil)
		return
	}

	e, err := h.svc.GetEvent(r.Context(), id)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	dtoEvent := dto.Event{
		ID:        e.GetID(),
		ProjectID: e.AggregateID(),
		Type:      e.Type(),
		CreatedBy: e.CreatedByID(),
		At:        e.OccurredAt(),
		Details:   e.String(),
	}

	switch ev := e.(type) {
	case events.ProjectRoleAssigned:
		dtoEvent.PersonID = &ev.PersonID
	case events.ProjectRoleUnassigned:
		dtoEvent.PersonID = &ev.PersonID
	case events.ProductAdded:
		dtoEvent.ProductID = &ev.ProductID
	case events.ProductRemoved:
		dtoEvent.ProductID = &ev.ProductID
	}

	_ = httputil.WriteJSON(w, http.StatusOK, dtoEvent)
}
