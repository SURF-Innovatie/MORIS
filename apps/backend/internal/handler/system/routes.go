package system

import (
	"github.com/go-chi/chi/v5"
)

// MountRoutes mounts all system API endpoints
func MountRoutes(r chi.Router, h *Handler) {
	r.Get("/health", h.Health)
	r.Get("/status", h.Status)
}
