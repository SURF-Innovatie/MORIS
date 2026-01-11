package eventpolicy

import (
	"github.com/go-chi/chi/v5"
)

// RegisterRoutes registers event policy routes
func (h *Handler) RegisterRoutes(r chi.Router) {
	// Policy CRUD routes
	r.Route("/policies", func(r chi.Router) {
		r.Get("/{id}", h.GetPolicy)
		r.Put("/{id}", h.UpdatePolicy)
		r.Delete("/{id}", h.DeletePolicy)
	})
}

// RegisterOrgRoutes registers org-scoped policy routes (to be mounted under /organisations)
func (h *Handler) RegisterOrgRoutes(r chi.Router) {
	r.Get("/{id}/policies", h.ListForOrgNode)
	r.Post("/{id}/policies", h.CreateForOrgNode)
}

// RegisterProjectRoutes registers project-scoped policy routes (to be mounted under /projects)
// Note: Policy creation/deletion for projects is managed via events (project.event_policy_added/removed)
// for proper RBAC control through project roles. Only read access is provided here.
func (h *Handler) RegisterProjectRoutes(r chi.Router) {
	r.Get("/{id}/policies", h.ListForProject)
}
