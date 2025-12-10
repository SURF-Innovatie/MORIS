package event

import (
	"encoding/json"
	"net/http"

	"github.com/SURF-Innovatie/MORIS/internal/api/eventdto"
	"github.com/SURF-Innovatie/MORIS/internal/api/persondto"
	"github.com/SURF-Innovatie/MORIS/internal/api/productdto"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
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

func toPersonDTO(p entities.Person) persondto.Response {
	return persondto.Response{
		ID:          p.Id,
		UserID:      p.UserID,
		Name:        p.Name,
		GivenName:   p.GivenName,
		FamilyName:  p.FamilyName,
		Email:       p.Email,
		AvatarUrl:   p.AvatarUrl,
		ORCiD:       p.ORCiD,
		Description: p.Description,
	}
}

func toProductDTO(p entities.Product) productdto.Response {
	return productdto.Response{
		ID:       p.Id,
		Name:     p.Name,
		Language: p.Language,
		Type:     p.Type,
		DOI:      p.DOI,
	}
}

// ApproveEvent godoc
// @Summary Approve an event
// @Description Approves a pending event
// @Tags events
// @Accept json
// @Produce json
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// RejectEvent godoc
// @Summary Reject an event
// @Description Rejects a pending event
// @Tags events
// @Accept json
// @Produce json
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// GetEvent godoc
// @Summary Get event details
// @Description Retrieves details for a specific event by ID
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "Event ID (UUID)"
// @Success 200 {object} eventdto.Event
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

	dtoEvent := eventdto.Event{
		ID:        e.GetID(),
		ProjectID: e.AggregateID(),
		Type:      e.Type(),
		Status:    e.GetStatus(),
		CreatedBy: e.CreatedByID(),
		At:        e.OccurredAt(),
		Details:   e.String(),
	}

	switch ev := e.(type) {
	case events.PersonAdded:
		p := toPersonDTO(ev.Person)
		dtoEvent.Person = &p
	case events.PersonRemoved:
		p := toPersonDTO(ev.Person)
		dtoEvent.Person = &p
	case events.ProductAdded:
		p := toProductDTO(ev.Product)
		dtoEvent.Product = &p
	case events.ProductRemoved:
		p := toProductDTO(ev.Product)
		dtoEvent.Product = &p
	}

	_ = httputil.WriteJSON(w, http.StatusOK, dtoEvent)
}
