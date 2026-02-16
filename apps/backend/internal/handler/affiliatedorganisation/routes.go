package affiliatedorganisation

import (
	"github.com/go-chi/chi/v5"
)

// MountRoutes mounts affiliated organisation routes to the router.
func MountRoutes(r chi.Router, h *Handler) {
	r.Route("/affiliated-organisations", func(r chi.Router) {
		// VAT lookup (must be before /{id} to avoid conflict)
		r.Get("/vat/lookup", h.LookupVAT)

		// CRUD
		r.Get("/", h.GetAll)
		r.Post("/", h.Create)
		r.Get("/{id}", h.Get)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}
