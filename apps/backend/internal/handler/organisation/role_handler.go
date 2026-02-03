package organisation

import (
	"net/http"

	"github.com/SURF-Innovatie/MORIS/internal/api/dto"
	organisationrbacsvc "github.com/SURF-Innovatie/MORIS/internal/app/organisation/rbac"
	organisationrolesvc "github.com/SURF-Innovatie/MORIS/internal/app/organisation/role"
	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
	"github.com/SURF-Innovatie/MORIS/internal/domain/organisation/rbac"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type RoleHandler struct {
	roleSvc organisationrolesvc.Service
	rbacSvc organisationrbacsvc.Service
}

func NewRoleHandler(roleSvc organisationrolesvc.Service, rbacSvc organisationrbacsvc.Service) *RoleHandler {
	return &RoleHandler{roleSvc: roleSvc, rbacSvc: rbacSvc}
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
func (h *RoleHandler) ListPermissions(w http.ResponseWriter, r *http.Request) {
	// Simple mapping for now - ideally this would be in the service or domain
	// But since it's just constants + display info, we can do it here or helper
	perms := lo.Map(rbac.Definitions, func(d rbac.PermissionDefinition, _ int) dto.PermissionDefinition {
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
func (h *RoleHandler) EnsureDefaultRoles(w http.ResponseWriter, r *http.Request) {
	if err := h.roleSvc.EnsureDefaultRoles(r.Context()); err != nil {
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
func (h *RoleHandler) ListRoles(w http.ResponseWriter, r *http.Request) {
	orgID, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid orgId", nil)
		return
	}

	roles, err := h.roleSvc.ListRoles(r.Context(), &orgID)
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
func (h *RoleHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	orgID, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid orgId", nil)
		return
	}

	if !requireOrgPerm(w, r, h.rbacSvc, orgID, rbac.PermissionManageOrganisationRoles) {
		return
	}

	var req dto.CreateRoleRequest
	if !httputil.ReadJSON(w, r, &req) {
		return
	}

	perms := lo.Map(req.Permissions, func(p string, _ int) rbac.Permission {
		return rbac.Permission(p)
	})

	newRole, err := h.roleSvc.CreateRole(r.Context(), orgID, req.Key, req.DisplayName, perms)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, transform.ToDTOItem[dto.OrganisationRoleResponse, *rbac.OrganisationRole](newRole))
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
func (h *RoleHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	roleID, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid id", nil)
		return
	}
	// TODO: We need the orgID to check permission.
	existingRole, err := h.roleSvc.GetRole(r.Context(), roleID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusNotFound, "role not found", nil)
		return
	}

	if !requireOrgPerm(w, r, h.rbacSvc, existingRole.OrganisationNodeID, rbac.PermissionManageOrganisationRoles) {
		return
	}

	var req dto.UpdateRoleRequest
	if !httputil.ReadJSON(w, r, &req) {
		return
	}

	perms := lo.Map(req.Permissions, func(p string, _ int) rbac.Permission {
		return rbac.Permission(p)
	})

	updatedRole, err := h.roleSvc.UpdateRole(r.Context(), roleID, req.DisplayName, perms)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, transform.ToDTOItem[dto.OrganisationRoleResponse, *rbac.OrganisationRole](updatedRole))
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
func (h *RoleHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	roleID, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid id", nil)
		return
	}

	existingRole, err := h.roleSvc.GetRole(r.Context(), roleID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusNotFound, "role not found", nil)
		return
	}

	if !requireOrgPerm(w, r, h.rbacSvc, existingRole.OrganisationNodeID, rbac.PermissionManageOrganisationRoles) {
		return
	}

	if err := h.roleSvc.DeleteRole(r.Context(), roleID); err != nil {
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
func (h *RoleHandler) CreateScope(w http.ResponseWriter, r *http.Request) {
	var req dto.OrganisationCreateScopeRequest
	if !httputil.ReadJSON(w, r, &req) {
		return
	}
	if req.RoleKey == "" {
		httputil.WriteError(w, r, http.StatusBadRequest, "roleKey is required", nil)
		return
	}

	if !requireOrgPerm(w, r, h.rbacSvc, req.RootNodeID, rbac.PermissionManageOrganisationRoles) {
		return
	}

	sc, err := h.roleSvc.CreateScope(r.Context(), req.RoleKey, req.RootNodeID)
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
func (h *RoleHandler) AddMembership(w http.ResponseWriter, r *http.Request) {
	var req dto.OrganisationAddMembershipRequest
	if !httputil.ReadJSON(w, r, &req) {
		return
	}

	// We need to know which organisation this scope belongs to, so we can check if caller manages members there.
	scope, err := h.roleSvc.GetScope(r.Context(), req.RoleScopeID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusNotFound, "scope not found", nil)
		return
	}

	if !requireOrgPerm(w, r, h.rbacSvc, scope.RootNodeID, rbac.PermissionManageMembers) {
		return
	}

	m, err := h.roleSvc.AddMembership(r.Context(), req.PersonID, req.RoleScopeID)
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
func (h *RoleHandler) RemoveMembership(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid id", nil)
		return
	}

	// Get membership to find scope -> org
	m, err := h.roleSvc.GetMembership(r.Context(), id)
	if err != nil {
		httputil.WriteError(w, r, http.StatusNotFound, "membership not found", nil)
		return
	}
	scope, err := h.roleSvc.GetScope(r.Context(), m.RoleScopeID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusNotFound, "scope not found", nil)
		return
	}

	if !requireOrgPerm(w, r, h.rbacSvc, scope.RootNodeID, rbac.PermissionManageMembers) {
		return
	}

	if err := h.roleSvc.RemoveMembership(r.Context(), id); err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteStatus(w)
}

func requireOrgPerm(w http.ResponseWriter, r *http.Request, rbac organisationrbacsvc.Service, nodeID uuid.UUID, perm rbac.Permission) bool {
	user, ok := httputil.GetUserFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return false
	}
	okPerm, err := rbac.HasPermission(r.Context(), user.Person.ID, nodeID, perm)
	if err != nil || !okPerm {
		httputil.WriteError(w, r, http.StatusForbidden, "forbidden", nil)
		return false
	}
	return true
}
