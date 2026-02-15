package organisation_rbac

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/identity"
	"github.com/SURF-Innovatie/MORIS/internal/domain/organisation"
	"github.com/SURF-Innovatie/MORIS/internal/domain/organisation/rbac"
	"github.com/google/uuid"
)

type repository interface {
	ListEffectiveMemberships(ctx context.Context, nodeID uuid.UUID) ([]EffectiveMembership, error)
	ListMyMemberships(ctx context.Context, personID uuid.UUID) ([]EffectiveMembership, error)
	GetMyPermissions(ctx context.Context, userID, nodeID uuid.UUID) ([]rbac.Permission, error)

	GetApprovalNode(ctx context.Context, nodeID uuid.UUID) (*organisation.OrganisationNode, error)
	HasAdminAccess(ctx context.Context, personID uuid.UUID, nodeID uuid.UUID) (bool, error)
	HasPermission(ctx context.Context, personID uuid.UUID, nodeID uuid.UUID, permission rbac.Permission) (bool, error)
}

type EffectiveMembership struct {
	MembershipID uuid.UUID
	PersonID     uuid.UUID

	RoleScopeID           uuid.UUID
	ScopeRootOrganisation *organisation.OrganisationNode

	RoleID         uuid.UUID
	RoleKey        string
	Permissions    []rbac.Permission
	HasAdminRights bool

	Person       identity.Person
	CustomFields map[string]any
}
