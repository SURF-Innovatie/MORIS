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

// ApproveEvent godoc
// @Summary Approve an event
// @Description Approves a pending event
// @Tags events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Event ID (UUID)"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "invalid event id"
// @Failure 500 {string} string "internal server error"
// @Router /events/{id}/approve [post]
func (h *Handler) ApproveEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid event id", nil)
		return
	}

	if err := h.svc.ApproveEvent(ctx, id); err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	httputil.WriteStatus(w)
}

// RejectEvent godoc
// @Summary Reject an event
// @Description Rejects a pending event
// @Tags events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Event ID (UUID)"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "invalid event id"
// @Failure 500 {string} string "internal server error"
// @Router /events/{id}/reject [post]
func (h *Handler) RejectEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid event id", nil)
		return
	}

	if err := h.svc.RejectEvent(ctx, id); err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	httputil.WriteStatus(w)
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
		Status:    e.GetStatus(),
		CreatedBy: e.CreatedByID(),
		At:        e.OccurredAt(),
		Details:   e.String(),
	}

	// Enrich DTO with related IDs
	if r, ok := e.(events.HasRelatedIDs); ok {
		ids := r.RelatedIDs()
		dtoEvent.PersonID = ids.PersonID
		dtoEvent.ProductID = ids.ProductID
	}

	_ = httputil.WriteJSON(w, http.StatusOK, dtoEvent)
}

// ListEventTypes godoc
// @Summary List all event types
// @Description Lists all event types and whether the current user is allowed to trigger them
// @Tags events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} []dto.EventTypeResponse
// @Failure 500 {string} string "internal server error"
// @Router /events/types [get]
func (h *Handler) ListEventTypes(w http.ResponseWriter, r *http.Request) {
	types, err := h.svc.GetEventTypes(r.Context())
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, types)
}
