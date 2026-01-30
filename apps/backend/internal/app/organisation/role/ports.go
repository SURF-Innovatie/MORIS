package role

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type repository interface {
	EnsureDefaultRoles(ctx context.Context) error

	ListRoles(ctx context.Context, orgID *uuid.UUID) ([]*entities.OrganisationRole, error)
	CreateRole(ctx context.Context, orgID uuid.UUID, key, displayName string, permissions []entities.Permission) (*entities.OrganisationRole, error)
	GetRole(ctx context.Context, roleID uuid.UUID) (*entities.OrganisationRole, error)
	UpdateRole(ctx context.Context, roleID uuid.UUID, displayName string, permissions []entities.Permission) (*entities.OrganisationRole, error)
	DeleteRole(ctx context.Context, roleID uuid.UUID) error

	CreateScope(ctx context.Context, roleKey string, rootNodeID uuid.UUID) (*entities.RoleScope, error)
	GetScope(ctx context.Context, id uuid.UUID) (*entities.RoleScope, error)

	AddMembership(ctx context.Context, personID uuid.UUID, roleScopeID uuid.UUID) (*entities.Membership, error)
	GetMembership(ctx context.Context, membershipID uuid.UUID) (*entities.Membership, error)
	RemoveMembership(ctx context.Context, membershipID uuid.UUID) error
}
