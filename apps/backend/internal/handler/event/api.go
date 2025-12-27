package event

import "github.com/go-chi/chi/v5"

func MountEventRoutes(r chi.Router, h *Handler) {
	r.Route("/events", func(r chi.Router) {
		r.Post("/{id}/approve", h.ApproveEvent)
		r.Post("/{id}/reject", h.RejectEvent)
		r.Get("/types", h.ListEventTypes)
		r.Get("/{id}", h.GetEvent)
	})
}
