package projectrole

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, key, name string, orgNodeID uuid.UUID) (*entities.ProjectRole, error)
	GetByKeyAndOrg(ctx context.Context, key string, orgNodeID uuid.UUID) (*entities.ProjectRole, error)
	Delete(ctx context.Context, id uuid.UUID, orgNodeID uuid.UUID) error
	ListByOrgIDs(ctx context.Context, orgIDs []uuid.UUID) ([]entities.ProjectRole, error)
	List(ctx context.Context) ([]entities.ProjectRole, error)
}
