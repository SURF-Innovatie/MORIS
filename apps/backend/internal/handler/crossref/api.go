package crossref

import "github.com/go-chi/chi/v5"

func MountCrossrefRoutes(r chi.Router, h *Handler) {
	r.Route("/crossref", func(r chi.Router) {
		r.Get("/works", h.GetWork)
	})
}
