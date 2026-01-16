package command

import "github.com/go-chi/chi/v5"

func MountProjectCommandRouter(r chi.Router, h *Handler) {
	r.Get("/{id}/events", h.ListAvailableEvents)
	r.Post("/{id}/events", h.ExecuteEvent)
}
