package person

import "github.com/go-chi/chi/v5"

func MountPersonRoutes(r chi.Router, h *Handler) {
	r.Route("/people", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Get("/{id}", h.Get)
		r.Put("/{id}", h.Update)
	})
}
