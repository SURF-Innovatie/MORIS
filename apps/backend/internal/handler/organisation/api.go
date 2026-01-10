package organisation

import "github.com/go-chi/chi/v5"

func MountOrganisationRoutes(r chi.Router, h *Handler, rbac *RBACHandler) {
	r.Route("/organisation-nodes", func(r chi.Router) {
		r.Get("/ror/search", h.SearchROR)
		r.Post("/", h.CreateRoot)

		r.Get("/roots", h.ListRoots)
		r.Get("/{id}", h.GetOrganisationNode)
		r.Patch("/{id}", h.UpdateOrganisationNode)

		r.Post("/{id}/children", h.CreateChild)
		r.Get("/{id}/children", h.ListChildren)

		// RBAC on nodes
		r.Get("/{id}/memberships/effective", rbac.ListEffectiveMemberships)
		r.Get("/{id}/approval-node", rbac.GetApprovalNode)

		// Project Roles
		r.Post("/{id}/roles", h.CreateProjectRole)
		r.Get("/{id}/roles", h.ListProjectRoles)
		r.Patch("/{id}/roles/{roleId}", h.UpdateProjectRole)
		r.Delete("/{id}/roles/{roleId}", h.DeleteProjectRole)

		// Custom Custom Fields
		r.Post("/{id}/custom-fields", h.CreateCustomField)
		r.Get("/{id}/custom-fields", h.ListCustomFields)
		r.Delete("/{id}/custom-fields/{fieldId}", h.DeleteCustomField)

		r.Put("/{id}/members/{personId}/custom-fields", h.UpdateMemberCustomFields)
	})

	r.Route("/organisation-roles", func(r chi.Router) {
		r.Get("/", rbac.ListRoles)
		r.Post("/ensure-defaults", rbac.EnsureDefaultRoles)
	})

	r.Route("/organisation-scopes", func(r chi.Router) {
		r.Post("/", rbac.CreateScope)
	})

	r.Route("/organisation-memberships", func(r chi.Router) {
		r.Post("/", rbac.AddMembership)
		r.Delete("/{id}", rbac.RemoveMembership)
		r.Get("/mine", rbac.ListMyMemberships)
	})
}
