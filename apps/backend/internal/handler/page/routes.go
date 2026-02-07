package page

import (
	"github.com/go-chi/chi/v5"
)

func Routes(h *Handler) func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/{slug}", h.Get)
		r.Post("/", h.Create)
		r.Put("/{id}", h.Update)
		r.Get("/projects/{projectId}", h.ListByProject)
		r.Get("/users/{userId}", h.ListByUser)
	}
}
