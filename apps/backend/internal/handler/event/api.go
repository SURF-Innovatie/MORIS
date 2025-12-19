package event

import "github.com/go-chi/chi/v5"

func MountEventRoutes(r chi.Router, h *Handler) {
	r.Route("/events", func(r chi.Router) {
		r.Get("/{id}", h.GetEvent)
	})
}
