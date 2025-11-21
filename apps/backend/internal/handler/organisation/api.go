package organisation

import "github.com/go-chi/chi/v5"

func MountOrganisationRoutes(r chi.Router, h *Handler) {
	r.Route("/organisations", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Get("/{id}", h.Get)
	})
}
