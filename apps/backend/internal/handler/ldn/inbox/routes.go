package inbox

import (
	"github.com/go-chi/chi/v5"
)

// Routes mounts the LDN inbox routes on the given router.
func Routes(r chi.Router, h *Handler) {
	r.Route("/ldn", func(r chi.Router) {
		r.Post("/inbox", h.Receive)
		r.Get("/inbox", h.List)
	})
}
