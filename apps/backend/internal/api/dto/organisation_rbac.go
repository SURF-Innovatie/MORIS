package dto

import (
	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/organisation"
	"github.com/google/uuid"
)

type OrganisationRoleResponse struct {
	ID             uuid.UUID `json:"id"`
	Key            string    `json:"key"`
	HasAdminRights bool      `json:"hasAdminRights"`
}

func (r OrganisationRoleResponse) FromEntity(e entities.OrganisationRole) OrganisationRoleResponse {
	return OrganisationRoleResponse{
		ID:             e.ID,
		Key:            e.Key,
		HasAdminRights: e.HasAdminRights,
	}
}

type OrganisationEnsureDefaultsResponse struct {
	Status string `json:"status"`
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

	RoleID         uuid.UUID `json:"roleId"`
	RoleKey        string    `json:"roleKey"`
	HasAdminRights bool      `json:"hasAdminRights"`
}

func (r OrganisationEffectiveMembershipResponse) FromEntity(e organisation.EffectiveMembership) OrganisationEffectiveMembershipResponse {
	return OrganisationEffectiveMembershipResponse{
		MembershipID:          e.MembershipID,
		Person:                transform.ToDTOItem[PersonResponse](e.Person),
		RoleScopeID:           e.RoleScopeID,
		ScopeRootOrganisation: transform.ToDTOItem[OrganisationResponse](*e.ScopeRootOrganisation),
		RoleID:                e.RoleID,
		RoleKey:               e.RoleKey,
		HasAdminRights:        e.HasAdminRights,
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
