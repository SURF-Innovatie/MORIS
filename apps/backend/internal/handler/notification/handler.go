package notification

import (
	"net/http"

	"github.com/SURF-Innovatie/MORIS/internal/api/dto"
	appnotif "github.com/SURF-Innovatie/MORIS/internal/app/notification"
	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
)

type Handler struct {
	svc appnotif.Service
}

func NewHandler(svc appnotif.Service) *Handler {
	return &Handler{svc: svc}
}

// ListMe godoc
// @Summary List notifications for current user
// @Tags notifications
// @Produce json
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

	notifs, err := h.svc.ListForUser(r.Context(), userCtx.User.ID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, transform.ToDTOs[dto.NotificationResponse](notifs))
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
