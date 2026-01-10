package system

import (
	"net/http"

	"github.com/SURF-Innovatie/MORIS/internal/api/dto"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// Status godoc
// @Summary Check API status
// @Description Returns the current status of the API
// @Tags status
// @Accept json
// @Produce json
// @Success 200 {object} httputil.StatusResponse
// @Router /status [get]
func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	_ = httputil.WriteStatus(w)
}

// Health godoc
// @Summary Health check
// @Description Returns the health status of the API
// @Tags status
// @Accept json
// @Produce json
// @Success 200 {object} httputil.StatusResponse
// @Router /health [get]
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	_ = httputil.WriteStatus(w)
}

// EventTypeInfo represents an event type with its metadata
type EventTypeInfo struct {
	Type         string `json:"type"`
	FriendlyName string `json:"friendlyName"`
}

// ListEventTypes godoc
// @Summary List all available event types
// @Description Returns all event types that can be configured for role permissions
// @Tags system
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} EventTypeInfo
// @Router /event-types [get]
func (h *Handler) ListEventTypes(w http.ResponseWriter, r *http.Request) {
	metas := events.GetAllMetas()
	out := make([]EventTypeInfo, 0, len(metas))
	for _, m := range metas {
		out = append(out, EventTypeInfo{
			Type:         m.Type,
			FriendlyName: m.FriendlyName,
		})
	}
	_ = httputil.WriteJSON(w, http.StatusOK, out)
}

// EventTypeInfoDTO explicitly uses dto package for swagger mapping
var _ = dto.AvailableEvent{}
