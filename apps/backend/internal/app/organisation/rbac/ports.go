package organisation_rbac

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type repository interface {
	ListEffectiveMemberships(ctx context.Context, nodeID uuid.UUID) ([]EffectiveMembership, error)
	ListMyMemberships(ctx context.Context, personID uuid.UUID) ([]EffectiveMembership, error)
	GetMyPermissions(ctx context.Context, userID, nodeID uuid.UUID) ([]entities.Permission, error)

	GetApprovalNode(ctx context.Context, nodeID uuid.UUID) (*entities.OrganisationNode, error)
	HasAdminAccess(ctx context.Context, personID uuid.UUID, nodeID uuid.UUID) (bool, error)
	HasPermission(ctx context.Context, personID uuid.UUID, nodeID uuid.UUID, permission entities.Permission) (bool, error)
}

type EffectiveMembership struct {
	MembershipID uuid.UUID
	PersonID     uuid.UUID

	RoleScopeID           uuid.UUID
	ScopeRootOrganisation *entities.OrganisationNode

	RoleID         uuid.UUID
	RoleKey        string
	Permissions    []entities.Permission
	HasAdminRights bool

	Person       entities.Person
	CustomFields map[string]interface{}
}
