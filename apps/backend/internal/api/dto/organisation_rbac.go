package dto

import (
	organisationrbac "github.com/SURF-Innovatie/MORIS/internal/app/organisation/rbac"
	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type OrganisationRoleResponse struct {
	ID          uuid.UUID `json:"id"`
	Key         string    `json:"key"`
	DisplayName string    `json:"displayName"`
	Permissions []string  `json:"permissions"`
}

func (r OrganisationRoleResponse) FromEntity(e *entities.OrganisationRole) OrganisationRoleResponse {
	perms := make([]string, len(e.Permissions))
	for i, p := range e.Permissions {
		perms[i] = string(p)
	}
	return OrganisationRoleResponse{
		ID:          e.ID,
		Key:         e.Key,
		DisplayName: e.DisplayName,
		Permissions: perms,
	}
}

type PermissionDefinition struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

type GetPermissionsResponse struct {
	Permissions []PermissionDefinition `json:"permissions"`
}

type OrganisationEnsureDefaultsResponse struct {
	Status string `json:"status"`
}

type CreateRoleRequest struct {
	Key         string   `json:"key"`
	DisplayName string   `json:"displayName"`
	Permissions []string `json:"permissions"`
}

type UpdateRoleRequest struct {
	DisplayName string   `json:"displayName"`
	Permissions []string `json:"permissions"`
}

type OrganisationCreateScopeRequest struct {
	RoleKey    string    `json:"roleKey"`
	RootNodeID uuid.UUID `json:"rootNodeId"`
}

type OrganisationRoleScopeResponse struct {
	ID         uuid.UUID `json:"id"`
	RoleID     uuid.UUID `json:"roleId"`
	RootNodeID uuid.UUID `json:"rootNodeId"`
}

func (r OrganisationRoleScopeResponse) FromEntity(e *entities.RoleScope) OrganisationRoleScopeResponse {
	return OrganisationRoleScopeResponse{
		ID:         e.ID,
		RoleID:     e.RoleID,
		RootNodeID: e.RootNodeID,
	}
}

type OrganisationAddMembershipRequest struct {
	PersonID    uuid.UUID `json:"personId"`
	RoleScopeID uuid.UUID `json:"roleScopeId"`
}

type OrganisationMembershipResponse struct {
	ID          uuid.UUID `json:"id"`
	PersonID    uuid.UUID `json:"personId"`
	RoleScopeID uuid.UUID `json:"roleScopeId"`
}

func (r OrganisationMembershipResponse) FromEntity(e *entities.Membership) OrganisationMembershipResponse {
	return OrganisationMembershipResponse{
		ID:          e.ID,
		PersonID:    e.PersonID,
		RoleScopeID: e.RoleScopeID,
	}
}

type OrganisationEffectiveMembershipResponse struct {
	MembershipID uuid.UUID      `json:"membershipId"`
	Person       PersonResponse `json:"person"`

	RoleScopeID           uuid.UUID            `json:"roleScopeId"`
	ScopeRootOrganisation OrganisationResponse `json:"scopeRootOrganisation"`

	RoleID         uuid.UUID      `json:"roleId"`
	RoleKey        string         `json:"roleKey"`
	Permissions    []string       `json:"permissions"`
	HasAdminRights bool           `json:"hasAdminRights"`
	CustomFields   map[string]any `json:"customFields"`
}

func (r OrganisationEffectiveMembershipResponse) FromEntity(e organisationrbac.EffectiveMembership) OrganisationEffectiveMembershipResponse {
	perms := make([]string, len(e.Permissions))
	for i, p := range e.Permissions {
		perms[i] = string(p)
	}
	return OrganisationEffectiveMembershipResponse{
		MembershipID:          e.MembershipID,
		Person:                transform.ToDTOItem[PersonResponse](e.Person),
		RoleScopeID:           e.RoleScopeID,
		ScopeRootOrganisation: transform.ToDTOItem[OrganisationResponse](*e.ScopeRootOrganisation),
		RoleID:                e.RoleID,
		RoleKey:               e.RoleKey,
		Permissions:           perms,
		HasAdminRights:        e.HasAdminRights,
		CustomFields:          e.CustomFields,
	}
}

type OrganisationApprovalNodeResponse struct {
	ApprovalNodeID uuid.UUID `json:"approvalNodeId"`
}

func (r OrganisationApprovalNodeResponse) FromEntity(e *entities.OrganisationNode) OrganisationApprovalNodeResponse {
	return OrganisationApprovalNodeResponse{
		ApprovalNodeID: e.ID,
	}
}
