package auth

import (
	coreauth "github.com/SURF-Innovatie/MORIS/internal/auth"
	"github.com/SURF-Innovatie/MORIS/internal/handler/middleware"
	"github.com/go-chi/chi/v5"
)

// MountRoutes mounts all auth API endpoints
func MountRoutes(r chi.Router, authSvc coreauth.Service, h *Handler) {

	r.Post("/login", h.Login)

	// SURFconext (OIDC) login
	r.Get("/auth/surfconext/url", h.GetSurfconextAuthURL)
	r.Post("/auth/surfconext/login", h.LoginWithSurfconext)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(authSvc))

		r.Get("/profile", h.Profile)

		// ORCID
		r.Get("/auth/orcid/url", h.GetORCIDAuthURL)
		r.Post("/auth/orcid/link", h.LinkORCID)
		r.Post("/auth/orcid/unlink", h.UnlinkORCID)
	})
}
