package event

import (
	"net/http"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/api/dto"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/queries"
	"github.com/SURF-Innovatie/MORIS/internal/app/user"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/SURF-Innovatie/MORIS/internal/event"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type Handler struct {
	svc      event.Service
	querySvc queries.Service
	userSvc  user.Service
	cli      *ent.Client
}

func NewHandler(svc event.Service, querySvc queries.Service, userSvc user.Service, cli *ent.Client) *Handler {
	return &Handler{svc: svc, querySvc: querySvc, userSvc: userSvc, cli: cli}
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

	detailed := events.DetailedEvent{Event: e}
	if hr, ok := e.(events.HasRelatedIDs); ok {
		ids := hr.RelatedIDs()
		if ids.PersonID != nil {
			people, _ := h.querySvc.GetPeopleByIDs(r.Context(), []uuid.UUID{*ids.PersonID})
			if p, ok := people[*ids.PersonID]; ok {
				detailed.Person = &p
			}
		}
		if ids.ProjectRoleID != nil {
			roles, _ := h.querySvc.GetProjectRolesByIDs(r.Context(), []uuid.UUID{*ids.ProjectRoleID})
			if role, ok := roles[*ids.ProjectRoleID]; ok {
				detailed.ProjectRole = &role
			}
		}
		if ids.ProductID != nil {
			products, _ := h.querySvc.GetProductsByIDs(r.Context(), []uuid.UUID{*ids.ProductID})
			if product, ok := products[*ids.ProductID]; ok {
				detailed.Product = &product
			}
		}
	}

	// Resolve Creator
	creatorID := e.CreatedByID()
	people, err := h.userSvc.GetPeopleByUserIDs(r.Context(), []uuid.UUID{creatorID})
	if err == nil {
		if p, ok := people[creatorID]; ok {
			detailed.Creator = &p
		}
	}

	var d dto.Event
	dtoEvent := d.FromDetailedEntity(detailed)

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

	eventTypes := lo.FilterMap(types, func(t events.EventMeta, _ int) (dto.EventTypeResponse, bool) {
		ev, err := events.Create(t.Type)
		if err != nil {
			return dto.EventTypeResponse{}, false
		}

		allowed := t.IsAllowed(r.Context(), ev, h.cli)

		return dto.EventTypeResponse{
			Type:         t.Type,
			FriendlyName: t.FriendlyName,
			Allowed:      allowed,
		}, true
	})

	_ = httputil.WriteJSON(w, http.StatusOK, eventTypes)
}
