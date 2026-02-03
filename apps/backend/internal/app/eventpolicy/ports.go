package eventpolicy

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/policy"
	"github.com/google/uuid"
)

// repository defines the persistence interface for event policies
type repository interface {
	Create(ctx context.Context, policy policy.EventPolicy) (*policy.EventPolicy, error)
	Update(ctx context.Context, id uuid.UUID, policy policy.EventPolicy) (*policy.EventPolicy, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*policy.EventPolicy, error)

	// ListForOrgNode returns policies for an org node.
	// If ancestorNodeIDs is provided, also returns policies from those ancestors (for inheritance).
	ListForOrgNode(ctx context.Context, orgNodeID uuid.UUID, ancestorNodeIDs []uuid.UUID) ([]policy.EventPolicy, error)

	// ListForProject returns policies directly attached to a project.
	ListForProject(ctx context.Context, projectID uuid.UUID) ([]policy.EventPolicy, error)
}

// RecipientResolver resolves recipient specifications to actual user IDs
type RecipientResolver interface {
	// ResolveUsers converts person IDs to user IDs (since policies store person IDs as "user IDs")
	ResolveUsers(ctx context.Context, personIDs []uuid.UUID) ([]uuid.UUID, error)

	// ResolveRole returns user IDs for users with a given project role
	ResolveRole(ctx context.Context, roleID uuid.UUID, projectID uuid.UUID) ([]uuid.UUID, error)

	// ResolveOrgRole returns user IDs for users with a given organisation role
	ResolveOrgRole(ctx context.Context, roleID uuid.UUID, orgNodeID uuid.UUID) ([]uuid.UUID, error)

	// ResolveDynamic returns user IDs for dynamic recipient types
	// dynType can be: "project_members", "project_owner", "org_admins"
	ResolveDynamic(ctx context.Context, dynType string, projectID uuid.UUID, orgNodeID uuid.UUID) ([]uuid.UUID, error)
}
