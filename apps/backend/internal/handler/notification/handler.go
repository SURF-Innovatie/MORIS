package notification

import (
	"encoding/json"
	"net/http"

	"github.com/SURF-Innovatie/MORIS/internal/auth"
	_ "github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/projectnotification"
)

type Handler struct {
	svc projectnotification.Service
}

func NewHandler(svc projectnotification.Service) *Handler {
	return &Handler{svc: svc}
}

// GetNotifications godoc
// @Summary Get notifications for the logged-in user
// @Description Retrieves a list of notifications for the authenticated user
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {array} entities.ProjectNotification
// @Failure 401 {string} string "unauthorized"
// @Failure 500 {string} string "internal server error"
// @Router /notifications [get]
func (h *Handler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.GetUserFromContext(r.Context())
	if !ok || user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	notifications, err := h.svc.ListForUser(r.Context(), user.User.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(notifications)
}
