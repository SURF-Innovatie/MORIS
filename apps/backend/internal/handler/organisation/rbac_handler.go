package organisation

import (
	"net/http"

	"github.com/SURF-Innovatie/MORIS/internal/api/dto"
	rbacsvc "github.com/SURF-Innovatie/MORIS/internal/app/organisation/rbac"
	"github.com/SURF-Innovatie/MORIS/internal/app/organisation/role"
	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
	"github.com/samber/lo"
)

type RBACHandler struct {
	rbac rbacsvc.Service
}

func NewRBACHandler(r rbacsvc.Service) *RBACHandler {
	return &RBACHandler{rbac: r}
}

// ListPermissions godoc
// @Summary List available permissions
// @Description Lists all available organisation permissions
// @Tags organisation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.GetPermissionsResponse
// @Failure 401 {string} string "unauthorized"
// @Failure 500 {string} string "internal server error"
// @Router /organisation-permissions [get]
func (h *RBACHandler) ListPermissions(w http.ResponseWriter, r *http.Request) {
	// Simple mapping for now - ideally this would be in the service or domain
	// But since it's just constants + display info, we can do it here or helper
	perms := lo.Map(role.Definitions, func(d role.PermissionDefinition, _ int) dto.PermissionDefinition {
		return dto.PermissionDefinition{
			Key:         string(d.Permission),
			Label:       d.Label,
			Description: d.Description,
		}
	})

	_ = httputil.WriteJSON(w, http.StatusOK, dto.GetPermissionsResponse{Permissions: perms})
}

// EnsureDefaultRoles godoc
// @Summary Ensure default organisation roles
// @Description Creates/updates default roles: admin, researcher, students
// @Tags organisation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.OrganisationEnsureDefaultsResponse
// @Failure 401 {string} string "unauthorized"
// @Failure 500 {string} string "internal server error"
// @Router /organisation-roles/ensure-defaults [post]
func (h *RBACHandler) EnsureDefaultRoles(w http.ResponseWriter, r *http.Request) {
	if err := h.rbac.EnsureDefaultRoles(r.Context()); err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	_ = httputil.WriteJSON(w, http.StatusOK, dto.OrganisationEnsureDefaultsResponse{Status: "ok"})
}

// ListRoles godoc
// @Summary List organisation roles for an organisation
// @Description Lists all organisation roles for a specific organisation
// @Tags organisation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Organisation ID"
// @Success 200 {array} dto.OrganisationRoleResponse
// @Failure 400 {string} string "invalid id"
// @Failure 401 {string} string "unauthorized"
// @Failure 500 {string} string "internal server error"
// @Router /organisation-nodes/{id}/organisation-roles [get]
func (h *RBACHandler) ListRoles(w http.ResponseWriter, r *http.Request) {
	orgID, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid orgId", nil)
		return
	}

	roles, err := h.rbac.ListRoles(r.Context(), &orgID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	out := transform.ToDTOs[dto.OrganisationRoleResponse](roles)

	_ = httputil.WriteJSON(w, http.StatusOK, out)
}

// CreateRole godoc
// @Summary Create a new organisation role
// @Description Creates a new role for an organisation
// @Tags organisation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Organisation ID"
// @Param request body dto.CreateRoleRequest true "Role details"
// @Success 200 {object} dto.OrganisationRoleResponse
// @Failure 400 {string} string "invalid id"
// @Failure 401 {string} string "unauthorized"
// @Failure 500 {string} string "internal server error"
// @Router /organisation-nodes/{id}/organisation-roles [post]
func (h *RBACHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	orgID, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid orgId", nil)
		return
	}

	user, ok := httputil.GetUserFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	if ok, err := h.rbac.HasPermission(r.Context(), user.Person.ID, orgID, role.PermissionManageOrganisationRoles); err != nil || !ok {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	var req dto.CreateRoleRequest
	if !httputil.ReadJSON(w, r, &req) {
		return
	}

	perms := lo.Map(req.Permissions, func(p string, _ int) role.Permission {
		return role.Permission(p)
	})

	newRole, err := h.rbac.CreateRole(r.Context(), orgID, req.Key, req.DisplayName, perms)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, transform.ToDTOItem[dto.OrganisationRoleResponse, *entities.OrganisationRole](newRole))
}

// UpdateRole godoc
// @Summary Update an organisation role
// @Description Updates an existing organisation role
// @Tags organisation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Role ID"
// @Param request body dto.UpdateRoleRequest true "Role details"
// @Success 200 {object} dto.OrganisationRoleResponse
// @Failure 400 {string} string "invalid id"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "not found"
// @Failure 500 {string} string "internal server error"
// @Router /organisation-roles/{id} [put]
func (h *RBACHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	roleID, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid id", nil)
		return
	}
	// TODO: We need the orgID to check permission.
	existingRole, err := h.rbac.GetRole(r.Context(), roleID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusNotFound, "role not found", nil)
		return
	}

	user, ok := httputil.GetUserFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	if ok, err := h.rbac.HasPermission(r.Context(), user.Person.ID, existingRole.OrganisationNodeID, role.PermissionManageOrganisationRoles); err != nil || !ok {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	var req dto.UpdateRoleRequest
	if !httputil.ReadJSON(w, r, &req) {
		return
	}

	perms := lo.Map(req.Permissions, func(p string, _ int) role.Permission {
		return role.Permission(p)
	})

	updatedRole, err := h.rbac.UpdateRole(r.Context(), roleID, req.DisplayName, perms)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, transform.ToDTOItem[dto.OrganisationRoleResponse, *entities.OrganisationRole](updatedRole))
}

// DeleteRole godoc
// @Summary Delete an organisation role
// @Description Deletes an organisation role
// @Tags organisation
// @Security BearerAuth
// @Param id path string true "Role ID"
// @Success 200
// @Failure 400 {string} string "invalid id"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "not found"
// @Failure 500 {string} string "internal server error"
// @Router /organisation-roles/{id} [delete]
func (h *RBACHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	roleID, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid id", nil)
		return
	}

	existingRole, err := h.rbac.GetRole(r.Context(), roleID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusNotFound, "role not found", nil)
		return
	}

	user, ok := httputil.GetUserFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	if ok, err := h.rbac.HasPermission(r.Context(), user.Person.ID, existingRole.OrganisationNodeID, role.PermissionManageOrganisationRoles); err != nil || !ok {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	if err := h.rbac.DeleteRole(r.Context(), roleID); err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteStatus(w)
}

// CreateScope godoc
// @Summary Create role scope
// @Description Creates a new role scope
// @Tags organisation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.OrganisationCreateScopeRequest true "Scope details"
// @Success 200 {object} dto.OrganisationRoleScopeResponse
// @Failure 401 {string} string "unauthorized"
// @Failure 500 {string} string "internal server error"
// @Router /organisation-scopes [post]
func (h *RBACHandler) CreateScope(w http.ResponseWriter, r *http.Request) {
	var req dto.OrganisationCreateScopeRequest
	if !httputil.ReadJSON(w, r, &req) {
		return
	}
	if req.RoleKey == "" {
		httputil.WriteError(w, r, http.StatusBadRequest, "roleKey is required", nil)
		return
	}

	user, ok := httputil.GetUserFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	if ok, err := h.rbac.HasPermission(r.Context(), user.Person.ID, req.RootNodeID, role.PermissionManageOrganisationRoles); err != nil || !ok {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	sc, err := h.rbac.CreateScope(r.Context(), req.RoleKey, req.RootNodeID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, transform.ToDTOItem[dto.OrganisationRoleScopeResponse](sc))
}

// AddMembership godoc
// @Summary Add membership
// @Description Adds a user membership to a scope
// @Tags organisation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.OrganisationAddMembershipRequest true "Membership details"
// @Success 200 {object} dto.OrganisationMembershipResponse
// @Failure 401 {string} string "unauthorized"
// @Failure 500 {string} string "internal server error"
// @Router /organisation-memberships [post]
func (h *RBACHandler) AddMembership(w http.ResponseWriter, r *http.Request) {
	var req dto.OrganisationAddMembershipRequest
	if !httputil.ReadJSON(w, r, &req) {
		return
	}

	// We need to know which organisation this scope belongs to, so we can check if caller manages members there.
	scope, err := h.rbac.GetScope(r.Context(), req.RoleScopeID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusNotFound, "scope not found", nil)
		return
	}

	user, ok := httputil.GetUserFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	if ok, err := h.rbac.HasPermission(r.Context(), user.Person.ID, scope.RootNodeID, role.PermissionManageMembers); err != nil || !ok {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	m, err := h.rbac.AddMembership(r.Context(), req.PersonID, req.RoleScopeID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, transform.ToDTOItem[dto.OrganisationMembershipResponse](m))
}

// RemoveMembership godoc
// @Summary Remove membership
// @Description Removes a user membership
// @Tags organisation
// @Security BearerAuth
// @Param id path string true "Membership ID"
// @Success 200
// @Failure 400 {string} string "invalid id"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "not found"
// @Failure 500 {string} string "internal server error"
// @Router /organisation-memberships/{id} [delete]
func (h *RBACHandler) RemoveMembership(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid id", nil)
		return
	}

	// Get membership to find scope -> org
	m, err := h.rbac.GetMembership(r.Context(), id)
	if err != nil {
		httputil.WriteError(w, r, http.StatusNotFound, "membership not found", nil)
		return
	}
	scope, err := h.rbac.GetScope(r.Context(), m.RoleScopeID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusNotFound, "scope not found", nil)
		return
	}

	user, ok := httputil.GetUserFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	if ok, err := h.rbac.HasPermission(r.Context(), user.Person.ID, scope.RootNodeID, role.PermissionManageMembers); err != nil || !ok {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	if err := h.rbac.RemoveMembership(r.Context(), id); err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteStatus(w)
}

// ListEffectiveMemberships godoc
// @Summary List effective memberships at a node
// @Description Lists all memberships whose scope root covers this node (scope root is an ancestor of the node)
// @Tags organisation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "OwningOrgNode node ID"
// @Success 200 {array} dto.OrganisationEffectiveMembershipResponse
// @Failure 400 {string} string "invalid id"
// @Failure 401 {string} string "unauthorized"
// @Failure 500 {string} string "internal server error"
// @Router /organisation-nodes/{id}/memberships/effective [get]
func (h *RBACHandler) ListEffectiveMemberships(w http.ResponseWriter, r *http.Request) {
	nodeID, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid id", nil)
		return
	}

	ms, err := h.rbac.ListEffectiveMemberships(r.Context(), nodeID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	out := transform.ToDTOs[dto.OrganisationEffectiveMembershipResponse](ms)

	_ = httputil.WriteJSON(w, http.StatusOK, out)
}

// GetApprovalNode godoc
// @Summary Get approval node for a node
// @Description Returns the node responsible for approvals by bubbling up to the nearest ancestor with an Admin scope rooted there that has members
// @Tags organisation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "OwningOrgNode node ID"
// @Success 200 {object} dto.OrganisationApprovalNodeResponse
// @Failure 400 {string} string "invalid id"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "not found"
// @Failure 500 {string} string "internal server error"
// @Router /organisation-nodes/{id}/approval-node [get]
func (h *RBACHandler) GetApprovalNode(w http.ResponseWriter, r *http.Request) {
	nodeID, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid id", nil)
		return
	}

	n, err := h.rbac.GetApprovalNode(r.Context(), nodeID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusNotFound, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, transform.ToDTOItem[dto.OrganisationApprovalNodeResponse](n))
}

// ListMyMemberships godoc
// @Summary List my memberships
// @Description Lists all memberships for the current user
// @Tags organisation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.OrganisationEffectiveMembershipResponse
// @Failure 401 {string} string "unauthorized"
// @Failure 500 {string} string "internal server error"
// @Router /organisation-memberships/mine [get]
func (h *RBACHandler) ListMyMemberships(w http.ResponseWriter, r *http.Request) {
	user, ok := httputil.GetUserFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	ms, err := h.rbac.ListMyMemberships(r.Context(), user.Person.ID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	out := transform.ToDTOs[dto.OrganisationEffectiveMembershipResponse](ms)

	_ = httputil.WriteJSON(w, http.StatusOK, out)
}

// GetMyPermissions godoc
// @Summary Get my effective permissions for a node
// @Description Returns the union of all permissions the current user has on the specified node (inherited/direct)
// @Tags organisation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Node ID"
// @Success 200 {object} []string
// @Failure 401 {string} string "unauthorized"
// @Failure 500 {string} string "internal server error"
// @Router /organisation-nodes/{id}/permissions/mine [get]
func (h *RBACHandler) GetMyPermissions(w http.ResponseWriter, r *http.Request) {
	nodeID, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid id", nil)
		return
	}

	user, ok := httputil.GetUserFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	perms, err := h.rbac.GetMyPermissions(r.Context(), user.Person.ID, nodeID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Convert to []string
	out := lo.Map(perms, func(p role.Permission, _ int) string {
		return string(p)
	})

	_ = httputil.WriteJSON(w, http.StatusOK, out)
}
