package nwo

import "github.com/go-chi/chi/v5"

// MountRoutes mounts NWO-related routes to the given router
func MountRoutes(r chi.Router, h *Handler) {
	r.Route("/nwo", func(r chi.Router) {
		r.Get("/projects", h.GetProjects)
		r.Get("/project", h.GetProject)
	})
}
