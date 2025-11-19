package project

import "github.com/go-chi/chi/v5"

func MountProjectRoutes(r chi.Router, h *Handler) {
	r.Route("/projects", func(r chi.Router) {
		r.Get("/", h.GetAllProjects)
		r.Post("/", h.StartProject)
		r.Get("/{id}", h.GetProject)
		r.Post("/{id}/people", h.AddPerson)
		r.Delete("/{id}/people/{personId}", h.RemovePerson)
	})
}
