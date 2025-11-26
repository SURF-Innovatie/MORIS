package custom

import (
	"github.com/SURF-Innovatie/MORIS/internal/auth"
	"github.com/go-chi/chi/v5"
)

// MountCustomHandlers mounts all custom API endpoints
func MountCustomHandlers(r chi.Router, h *Handler) {
	r.Get("/health", h.Health)
	r.Get("/status", h.Status)
	r.Post("/register", h.Register)
	r.Post("/login", h.Login)

	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware)

		r.Get("/profile", h.Profile)
		r.Get("/users/count", h.TotalUserCount)

		// ORCID
		r.Get("/auth/orcid/url", h.GetORCIDAuthURL)
		r.Post("/auth/orcid/link", h.LinkORCID)
		r.Post("/auth/orcid/unlink", h.UnlinkORCID)

		r.Group(func(r chi.Router) {
			r.Use(auth.RequireRoleMiddleware("admin"))
			r.Get("/admin/users/list", h.AdminUserList)
		})
	})
}
