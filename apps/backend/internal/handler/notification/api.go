package notification

import "github.com/go-chi/chi/v5"

func MountNotificationRoutes(r chi.Router, h *Handler) {
	r.Route("/notifications", func(r chi.Router) {
		r.Get("/", h.GetNotifications)
	})
}
