package project

import "github.com/go-chi/chi/v5"

func MountProjectRoutes(r chi.Router, h *Handler) {
	r.Get("/", h.GetAllProjects)
	r.Get("/slug-check", h.CheckSlugAvailability)
	r.Get("/{id}", h.GetProject)
	r.Get("/{id}/changelog", h.GetChangelog)
	r.Get("/{id}/pending-events", h.GetPendingEvents)
	r.Get("/{id}/allowed-events", h.GetAllowedEvents)
	r.Get("/{id}/custom-fields", h.ListAvailableCustomFields)
}
