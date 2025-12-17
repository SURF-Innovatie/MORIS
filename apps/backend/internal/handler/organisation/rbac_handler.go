package organisation

import (
	"net/http"

	"github.com/SURF-Innovatie/MORIS/internal/api/dto"
	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
	organisationsvc "github.com/SURF-Innovatie/MORIS/internal/organisation"
)

type RBACHandler struct {
	rbac organisationsvc.RBACService
}

func NewRBACHandler(r organisationsvc.RBACService) *RBACHandler {
	return &RBACHandler{rbac: r}
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
// @Summary List organisation roles
// @Description Lists all organisation roles
// @Tags organisation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.OrganisationRoleResponse
// @Failure 401 {string} string "unauthorized"
// @Failure 500 {string} string "internal server error"
// @Router /organisation-roles [get]
func (h *RBACHandler) ListRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := h.rbac.ListRoles(r.Context())
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	out := transform.ToDTOs[dto.OrganisationRoleResponse](roles)

	_ = httputil.WriteJSON(w, http.StatusOK, out)
}

// CreateScope godoc
// @Summary Create a role scope
// @Description Creates a role scope for a role key and a root node (scope applies to root node and its subtree)
// @Tags organisation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.OrganisationCreateScopeRequest true "Create scope request"
// @Success 200 {object} dto.OrganisationRoleScopeResponse
// @Failure 400 {string} string "invalid request body / roleKey is required"
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

	sc, err := h.rbac.CreateScope(r.Context(), req.RoleKey, req.RootNodeID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, transform.ToDTOItem[dto.OrganisationRoleScopeResponse](sc))
}

// AddMembership godoc
// @Summary Add a membership
// @Description Adds a person to a role scope
// @Tags organisation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.OrganisationAddMembershipRequest true "Add membership request"
// @Success 200 {object} dto.OrganisationMembershipResponse
// @Failure 400 {string} string "invalid request body"
// @Failure 401 {string} string "unauthorized"
// @Failure 500 {string} string "internal server error"
// @Router /organisation-memberships [post]
func (h *RBACHandler) AddMembership(w http.ResponseWriter, r *http.Request) {
	var req dto.OrganisationAddMembershipRequest
	if !httputil.ReadJSON(w, r, &req) {
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
// @Summary Remove a membership
// @Description Removes a membership by ID
// @Tags organisation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Membership ID"
// @Success 200 {object} httputil.StatusResponse
// @Failure 400 {string} string "invalid id"
// @Failure 401 {string} string "unauthorized"
// @Failure 500 {string} string "internal server error"
// @Router /organisation-memberships/{id} [delete]
func (h *RBACHandler) RemoveMembership(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid id", nil)
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
