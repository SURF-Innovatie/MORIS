package adapter

import (
	"github.com/go-chi/chi/v5"
)

func MountRoutes(r chi.Router, h *Handler) {
	r.Get("/adapters", h.ListAdapters)
	r.Post("/projects/{id}/export/{sink}", h.ExportProject)
}
