package projecthandler

import "github.com/go-chi/chi/v5"

func MountProjectRoutes(r chi.Router, h *Handler) {
	r.Route("/projects", func(r chi.Router) {
		r.Get("/", h.GetAllProjects)
		r.Get("/{id}", h.GetProject)
		r.Post("/{id}/person/add", h.AddPerson)
		r.Post("/{id}/person/remove", h.RemovePerson)
	})
}
