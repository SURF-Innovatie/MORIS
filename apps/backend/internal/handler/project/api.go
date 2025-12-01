package project

import "github.com/go-chi/chi/v5"

func MountProjectRoutes(r chi.Router, h *Handler) {
	r.Route("/projects", func(r chi.Router) {
		r.Get("/", h.GetAllProjects)
		r.Post("/", h.StartProject)
		r.Get("/{id}", h.GetProject)
		r.Put("/{id}", h.UpdateProject)
		r.Post("/{id}/people/{personId}", h.AddPerson)
		r.Delete("/{id}/people/{personId}", h.RemovePerson)
		r.Post("/{id}/products/{productID}", h.AddProduct)
		r.Delete("/{id}/products/{productID}", h.RemoveProduct)
		r.Get("/{id}/changelog", h.GetChangelog)
		r.Get("/{id}/pending-events", h.GetPendingEvents)
	})
}

func MountEventRoutes(r chi.Router, h *Handler) {
	r.Route("/events", func(r chi.Router) {
		r.Post("/{id}/approve", h.ApproveEvent)
		r.Post("/{id}/reject", h.RejectEvent)
	})
}
