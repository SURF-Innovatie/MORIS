package organisation_rbac

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type Repository interface {
	EnsureDefaultRoles(ctx context.Context) error

	ListRoles(ctx context.Context) ([]entities.OrganisationRole, error)

	CreateScope(ctx context.Context, roleKey string, rootNodeID uuid.UUID) (*entities.RoleScope, error)
	GetScope(ctx context.Context, id uuid.UUID) (*entities.RoleScope, error)

	AddMembership(ctx context.Context, personID uuid.UUID, roleScopeID uuid.UUID) (*entities.Membership, error)
	RemoveMembership(ctx context.Context, membershipID uuid.UUID) error

	ListEffectiveMemberships(ctx context.Context, nodeID uuid.UUID) ([]EffectiveMembership, error)
	ListMyMemberships(ctx context.Context, personID uuid.UUID) ([]EffectiveMembership, error)

	GetApprovalNode(ctx context.Context, nodeID uuid.UUID) (*entities.OrganisationNode, error)
	HasAdminAccess(ctx context.Context, personID uuid.UUID, nodeID uuid.UUID) (bool, error)
}

type EffectiveMembership struct {
	MembershipID uuid.UUID
	PersonID     uuid.UUID

	RoleScopeID           uuid.UUID
	ScopeRootOrganisation *entities.OrganisationNode

	RoleID         uuid.UUID
	RoleKey        string
	HasAdminRights bool

	Person entities.Person
}
