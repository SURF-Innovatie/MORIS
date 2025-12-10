package notification

import (
	"net/http"

	"github.com/SURF-Innovatie/MORIS/internal/api/notificationdto"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
	"github.com/SURF-Innovatie/MORIS/internal/notification"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	svc notification.Service
}

func NewHandler(svc notification.Service) *Handler {
	return &Handler{svc: svc}
}

// GetNotifications godoc
// @Summary Get notifications for the logged-in user
// @Description Retrieves a list of notifications for the authenticated user
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {array} notificationdto.Response
// @Failure 401 {string} string "unauthorized"
// @Failure 500 {string} string "internal server error"
// @Router /notifications [get]
func (h *Handler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	user, ok := httputil.GetUserFromContext(r.Context())
	if !ok || user == nil {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	notifications, err := h.svc.ListForUser(r.Context(), user.User.ID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		logrus.Errorf("failed to list notifications for user %s: %v", user.User.ID, err)
		return
	}

	var dtos []notificationdto.Response
	for _, n := range notifications {
		dtos = append(dtos, notificationdto.FromEntity(n))
	}

	_ = httputil.WriteJSON(w, http.StatusOK, dtos)
}

// MarkNotificationAsRead godoc
// @Summary Mark a notification as read
// @Description Marks a specific notification as read
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path string true "Notification ID"
// @Success 200 {string} string "ok"
// @Failure 401 {string} string "unauthorized"
// @Failure 500 {string} string "internal server error"
// @Router /notifications/{id}/read [put]
func (h *Handler) MarkNotificationAsRead(w http.ResponseWriter, r *http.Request) {
	user, ok := httputil.GetUserFromContext(r.Context())
	if !ok || user == nil {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid id", nil)
		return
	}

	notification, err := h.svc.Get(r.Context(), id)
	if err != nil {
		httputil.WriteError(w, r, http.StatusNotFound, err.Error(), nil)
		return
	}

	if notification.User.ID != user.User.ID {
		httputil.WriteError(w, r, http.StatusForbidden, "forbidden", nil)
		return
	}

	if notification.Read {
		w.WriteHeader(http.StatusOK)
		return
	}

	err = h.svc.MarkAsRead(r.Context(), id)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	w.WriteHeader(http.StatusOK)
}
