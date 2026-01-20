package portfolio

import "github.com/go-chi/chi/v5"

func MountPortfolioRoutes(r chi.Router, h *Handler) {
	r.Route("/portfolio", func(r chi.Router) {
		r.Get("/me", h.GetMe)
		r.Put("/me", h.UpdateMe)
	})
}
