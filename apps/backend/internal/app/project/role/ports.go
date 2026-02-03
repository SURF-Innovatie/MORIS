package role

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/project/role"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, key, name string, orgNodeID uuid.UUID) (*role.ProjectRole, error)
	CreateWithEventTypes(ctx context.Context, key, name string, orgNodeID uuid.UUID, allowedEventTypes []string) (*role.ProjectRole, error)
	CreateOrRestore(ctx context.Context, key, name string, orgNodeID uuid.UUID) (*role.ProjectRole, error)
	GetByKeyAndOrg(ctx context.Context, key string, orgNodeID uuid.UUID) (*role.ProjectRole, error)
	GetByID(ctx context.Context, id uuid.UUID) (*role.ProjectRole, error)
	Exists(ctx context.Context, key string, orgNodeID uuid.UUID) (bool, error)
	Unarchive(ctx context.Context, key string, orgNodeID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID, orgNodeID uuid.UUID) error
	ListByOrgIDs(ctx context.Context, orgIDs []uuid.UUID) ([]role.ProjectRole, error)
	List(ctx context.Context) ([]role.ProjectRole, error)
	UpdateAllowedEventTypes(ctx context.Context, id uuid.UUID, eventTypes []string) (*role.ProjectRole, error)
}
