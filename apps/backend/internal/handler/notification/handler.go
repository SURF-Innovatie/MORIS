package notification

import (
	"net/http"
	"strings"

	"github.com/SURF-Innovatie/MORIS/internal/api/dto"
	appnotif "github.com/SURF-Innovatie/MORIS/internal/app/notification"
	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
	"github.com/google/uuid"
)

type Handler struct {
	svc appnotif.Service
}

func NewHandler(svc appnotif.Service) *Handler {
	return &Handler{svc: svc}
}

// ListMe godoc
// @Summary List notifications for current user
// @Description Returns notifications in JSON format by default, or AS2 JSON-LD if Accept header is application/ld+json
// @Tags notifications
// @Produce json
// @Produce application/ld+json
// @Security BearerAuth
// @Success 200 {array} dto.NotificationResponse
// @Failure 401 {string} string "unauthorized"
// @Failure 500 {string} string "internal server error"
// @Router /notifications/me [get]
func (h *Handler) ListMe(w http.ResponseWriter, r *http.Request) {
	userCtx, ok := httputil.GetUserFromContext(r.Context())
	if !ok || userCtx == nil {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	// Content negotiation: check Accept header
	accept := r.Header.Get("Accept")
	if strings.Contains(accept, "application/ld+json") || strings.Contains(accept, "application/activity+json") {
		h.listMeAS2(w, r, userCtx.User.ID)
		return
	}

	notifs, err := h.svc.ListForUser(r.Context(), userCtx.User.ID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, transform.ToDTOs[dto.NotificationResponse](notifs))
}

// listMeAS2 returns notifications in AS2 JSON-LD format.
func (h *Handler) listMeAS2(w http.ResponseWriter, r *http.Request, userID any) {
	uid, ok := userID.(uuid.UUID)
	if !ok {
		httputil.WriteError(w, r, http.StatusInternalServerError, "invalid user id type", nil)
		return
	}

	activities, err := h.svc.ListForUserAsActivities(r.Context(), uid)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	w.Header().Set("Content-Type", "application/ld+json")
	_ = httputil.WriteJSON(w, http.StatusOK, activities)
}

// MarkAsRead godoc
// @Summary Mark notification as read
// @Tags notifications
// @Produce json
// @Security BearerAuth
// @Param id path string true "Notification ID (UUID)"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "invalid id"
// @Failure 500 {string} string "internal server error"
// @Router /notifications/{id}/read [post]
func (h *Handler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid id", nil)
		return
	}

	if err := h.svc.MarkAsRead(r.Context(), id); err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	w.WriteHeader(http.StatusOK)
}
