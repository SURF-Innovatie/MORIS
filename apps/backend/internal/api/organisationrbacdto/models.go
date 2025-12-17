package organisationrbacdto

import (
	"github.com/SURF-Innovatie/MORIS/internal/api/organisationdto"
	"github.com/SURF-Innovatie/MORIS/internal/api/userdto"
	"github.com/google/uuid"
)

type RoleResponse struct {
	ID             uuid.UUID `json:"id"`
	Key            string    `json:"key"`
	HasAdminRights bool      `json:"hasAdminRights"`
}

type EnsureDefaultsResponse struct {
	Status string `json:"status"`
}

type CreateScopeRequest struct {
	RoleKey    string    `json:"roleKey"`
	RootNodeID uuid.UUID `json:"rootNodeId"`
}

type RoleScopeResponse struct {
	ID         uuid.UUID `json:"id"`
	RoleID     uuid.UUID `json:"roleId"`
	RootNodeID uuid.UUID `json:"rootNodeId"`
}

type AddMembershipRequest struct {
	PersonID    uuid.UUID `json:"personId"`
	RoleScopeID uuid.UUID `json:"roleScopeId"`
}

type MembershipResponse struct {
	ID          uuid.UUID `json:"id"`
	PersonID    uuid.UUID `json:"personId"`
	RoleScopeID uuid.UUID `json:"roleScopeId"`
}

type EffectiveMembershipResponse struct {
	MembershipID uuid.UUID              `json:"membershipId"`
	Person       userdto.PersonResponse `json:"person"`

	RoleScopeID           uuid.UUID                `json:"roleScopeId"`
	ScopeRootOrganisation organisationdto.Response `json:"scopeRootOrganisation"`

	RoleID         uuid.UUID `json:"roleId"`
	RoleKey        string    `json:"roleKey"`
	HasAdminRights bool      `json:"hasAdminRights"`
}

type ApprovalNodeResponse struct {
	ApprovalNodeID uuid.UUID `json:"approvalNodeId"`
}
