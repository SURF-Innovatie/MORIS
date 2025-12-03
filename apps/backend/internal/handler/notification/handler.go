package notification

import (
	"encoding/json"
	"net/http"

	"github.com/SURF-Innovatie/MORIS/internal/api/notificationdto"
	"github.com/SURF-Innovatie/MORIS/internal/handler/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	_ "github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/notification"
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
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	notifications, err := h.svc.ListForUser(r.Context(), user.User.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logrus.Errorf("failed to list notifications for user %s: %v", user.User.ID, err)
		return
	}

	var dtos []notificationdto.Response
	for _, n := range notifications {
		projectId := uuid.Nil
		eventId := uuid.Nil
		if n.Event != nil {
			projectId = n.Event.ProjectID
			eventId = n.Event.ID
		}
		dtos = append(dtos, notificationdto.Response{
			ID:        n.Id,
			Message:   n.Message,
			Type:      n.Type,
			Read:      n.Read,
			ProjectID: projectId,
			EventID:   eventId,
			SentAt:    n.SentAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(dtos)
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
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	notification, err := h.svc.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if notification.User.ID != user.User.ID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	if notification.Read {
		w.WriteHeader(http.StatusOK)
		return
	}

	err = h.svc.MarkAsRead(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
