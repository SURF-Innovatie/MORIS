package user

import "github.com/go-chi/chi/v5"

func MountUserRoutes(r chi.Router, h *Handler) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/", h.CreateUser)
		r.Get("/{id}", h.GetUser)
		r.Put("/{id}", h.UpdateUser)
		r.Delete("/{id}", h.DeleteUser)
	})
}
