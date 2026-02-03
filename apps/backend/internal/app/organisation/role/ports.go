package role

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/organisation/rbac"
	"github.com/google/uuid"
)

type repository interface {
	EnsureDefaultRoles(ctx context.Context) error

	ListRoles(ctx context.Context, orgID *uuid.UUID) ([]*rbac.OrganisationRole, error)
	CreateRole(ctx context.Context, orgID uuid.UUID, key, displayName string, permissions []rbac.Permission) (*rbac.OrganisationRole, error)
	GetRole(ctx context.Context, roleID uuid.UUID) (*rbac.OrganisationRole, error)
	UpdateRole(ctx context.Context, roleID uuid.UUID, displayName string, permissions []rbac.Permission) (*rbac.OrganisationRole, error)
	DeleteRole(ctx context.Context, roleID uuid.UUID) error

	CreateScope(ctx context.Context, roleKey string, rootNodeID uuid.UUID) (*rbac.RoleScope, error)
	GetScope(ctx context.Context, id uuid.UUID) (*rbac.RoleScope, error)

	AddMembership(ctx context.Context, personID uuid.UUID, roleScopeID uuid.UUID) (*rbac.Membership, error)
	GetMembership(ctx context.Context, membershipID uuid.UUID) (*rbac.Membership, error)
	RemoveMembership(ctx context.Context, membershipID uuid.UUID) error
}
