package affiliatedorganisation

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

// Repository defines the persistence interface for AffiliatedOrganisation.
type Repository interface {
	Get(ctx context.Context, id uuid.UUID) (*entities.AffiliatedOrganisation, error)
	List(ctx context.Context) ([]*entities.AffiliatedOrganisation, error)
	Create(ctx context.Context, org entities.AffiliatedOrganisation) (*entities.AffiliatedOrganisation, error)
	Update(ctx context.Context, id uuid.UUID, org entities.AffiliatedOrganisation) (*entities.AffiliatedOrganisation, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]entities.AffiliatedOrganisation, error)
}
