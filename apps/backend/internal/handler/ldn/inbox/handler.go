// Package inbox provides the LDN Inbox handler for receiving notifications.
// Implements W3C LDN Receiver: https://www.w3.org/TR/ldn/#receiver
package inbox

import (
	"encoding/json"
	"net/http"

	appnotification "github.com/SURF-Innovatie/MORIS/internal/app/notification"
	"github.com/SURF-Innovatie/MORIS/internal/domain/ldn"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
)

// Handler handles LDN Inbox HTTP endpoints.
type Handler struct {
	svc appnotification.Service
}

// NewHandler creates a new LDN Inbox Handler.
func NewHandler(svc appnotification.Service) *Handler {
	return &Handler{svc: svc}
}

// Receive handles POST /ldn/inbox to receive notifications from external services.
// @Summary Receive LDN notification
// @Description Receives an LDN notification from an external service. Requires authentication.
// @Tags ldn
// @Accept application/ld+json
// @Produce json
// @Security BearerAuth
// @Param body body ldn.Activity true "AS2 Activity payload"
// @Success 201 {object} map[string]string "Notification received"
// @Failure 400 {string} string "Invalid request body"
// @Failure 401 {string} string "Unauthorized"
// @Failure 415 {string} string "Unsupported Media Type"
// @Failure 422 {string} string "Invalid activity payload"
// @Router /ldn/inbox [post]
func (h *Handler) Receive(w http.ResponseWriter, r *http.Request) {
	// Validate Content-Type per LDN spec
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/ld+json" && contentType != "application/activity+json" {
		httputil.WriteError(w, r, http.StatusUnsupportedMediaType,
			"Content-Type must be application/ld+json or application/activity+json", nil)
		return
	}

	// Require authentication
	userCtx, ok := httputil.GetUserFromContext(r.Context())
	if !ok || userCtx == nil {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	// Parse AS2 Activity payload
	var activity ldn.Activity
	if err := json.NewDecoder(r.Body).Decode(&activity); err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid JSON payload", err)
		return
	}

	// Validate required COAR Notify properties
	if err := activity.Validate(); err != nil {
		httputil.WriteError(w, r, http.StatusUnprocessableEntity, err.Error(), nil)
		return
	}

	// Process and store the activity using the notification service
	_, err := h.svc.CreateFromActivity(r.Context(), &activity, userCtx.User.ID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, "failed to process notification", err)
		return
	}

	// Return 201 with Location header per LDN spec
	w.Header().Set("Location", "/api/notifications/"+activity.ID)
	w.WriteHeader(http.StatusCreated)
	_ = httputil.WriteJSON(w, http.StatusCreated, map[string]string{
		"status":  "received",
		"id":      activity.ID,
		"message": "Notification received and processed",
	})
}

// List handles GET /ldn/inbox to list received notifications (optional per LDN spec).
// @Summary List LDN inbox notifications
// @Description Lists notifications received via LDN inbox
// @Tags ldn
// @Produce application/ld+json
// @Security BearerAuth
// @Success 200 {array} ldn.Activity "List of activities"
// @Failure 401 {string} string "Unauthorized"
// @Router /ldn/inbox [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	// Require authentication
	userCtx, ok := httputil.GetUserFromContext(r.Context())
	if !ok || userCtx == nil {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	// For now, return empty list - will be implemented with service
	w.Header().Set("Content-Type", "application/ld+json")
	_ = httputil.WriteJSON(w, http.StatusOK, []ldn.Activity{})
}
