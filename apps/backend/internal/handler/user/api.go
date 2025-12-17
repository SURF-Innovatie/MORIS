package user

import (
	"github.com/SURF-Innovatie/MORIS/internal/handler/middleware"
	"github.com/go-chi/chi/v5"
)

func MountUserRoutes(r chi.Router, h *Handler) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/", h.CreateUser)
		r.Get("/{id}", h.GetUser)
		r.Put("/{id}", h.UpdateUser)

		r.Delete("/{id}", h.DeleteUser)
		r.Get("/{id}/events/approved", h.GetApprovedEvents)
		r.Get("/search", h.SearchUsers)
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireSysAdminMiddleware())
		r.Get("/admin/users/list", h.ListUsers)
		r.Post("/admin/users/{id}/toggle-active", h.ToggleActive)
	})
}
