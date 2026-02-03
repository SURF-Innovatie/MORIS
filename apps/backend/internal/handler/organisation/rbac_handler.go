package organisation

import (
	"net/http"

	"github.com/SURF-Innovatie/MORIS/internal/api/dto"
	rbacsvc "github.com/SURF-Innovatie/MORIS/internal/app/organisation/rbac"
	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
	"github.com/SURF-Innovatie/MORIS/internal/domain/organisation/rbac"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
	"github.com/samber/lo"
)

type RBACHandler struct {
	rbac rbacsvc.Service
}

func NewRBACHandler(r rbacsvc.Service) *RBACHandler {
	return &RBACHandler{rbac: r}
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
	out := lo.Map(perms, func(p rbac.Permission, _ int) string {
		return string(p)
	})

	_ = httputil.WriteJSON(w, http.StatusOK, out)
}
