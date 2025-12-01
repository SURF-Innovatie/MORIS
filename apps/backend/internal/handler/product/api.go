package product

import "github.com/go-chi/chi/v5"

func MountProductRoutes(r chi.Router, h *Handler) {
	r.Route("/products", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Get("/", h.GetAll)
		r.Get("/me", h.GetMe)
		r.Get("/{id}", h.Get)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}
