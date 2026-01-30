package organisation_rbac

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/app/organisation/role"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type Repository interface {
	EnsureDefaultRoles(ctx context.Context) error

	// Roles
	ListRoles(ctx context.Context, orgID *uuid.UUID) ([]*entities.OrganisationRole, error)
	CreateRole(ctx context.Context, orgID uuid.UUID, key, displayName string, permissions []role.Permission) (*entities.OrganisationRole, error)
	GetRole(ctx context.Context, roleID uuid.UUID) (*entities.OrganisationRole, error) // Added GetRole
	UpdateRole(ctx context.Context, roleID uuid.UUID, displayName string, permissions []role.Permission) (*entities.OrganisationRole, error)
	DeleteRole(ctx context.Context, roleID uuid.UUID) error

	CreateScope(ctx context.Context, roleKey string, rootNodeID uuid.UUID) (*entities.RoleScope, error)
	GetScope(ctx context.Context, id uuid.UUID) (*entities.RoleScope, error)

	AddMembership(ctx context.Context, personID uuid.UUID, roleScopeID uuid.UUID) (*entities.Membership, error)
	GetMembership(ctx context.Context, membershipID uuid.UUID) (*entities.Membership, error)
	RemoveMembership(ctx context.Context, membershipID uuid.UUID) error

	ListEffectiveMemberships(ctx context.Context, nodeID uuid.UUID) ([]EffectiveMembership, error)
	ListMyMemberships(ctx context.Context, personID uuid.UUID) ([]EffectiveMembership, error)
	GetMyPermissions(ctx context.Context, userID, nodeID uuid.UUID) ([]role.Permission, error)

	GetApprovalNode(ctx context.Context, nodeID uuid.UUID) (*entities.OrganisationNode, error)
	HasAdminAccess(ctx context.Context, personID uuid.UUID, nodeID uuid.UUID) (bool, error)
	HasPermission(ctx context.Context, personID uuid.UUID, nodeID uuid.UUID, permission role.Permission) (bool, error)

	AncestorIDs(ctx context.Context, nodeID uuid.UUID) ([]uuid.UUID, error)
	IsAncestor(ctx context.Context, ancestorID, descendantID uuid.UUID) (bool, error)
}

type EffectiveMembership struct {
	MembershipID uuid.UUID
	PersonID     uuid.UUID

	RoleScopeID           uuid.UUID
	ScopeRootOrganisation *entities.OrganisationNode

	RoleID         uuid.UUID
	RoleKey        string
	Permissions    []role.Permission
	HasAdminRights bool

	Person       entities.Person
	CustomFields map[string]interface{}
}
