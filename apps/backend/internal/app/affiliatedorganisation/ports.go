package affiliatedorganisation

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/affiliatedorganisation"
	"github.com/google/uuid"
)

// Repository defines the persistence interface for AffiliatedOrganisation.
type Repository interface {
	Get(ctx context.Context, id uuid.UUID) (*affiliatedorganisation.AffiliatedOrganisation, error)
	List(ctx context.Context) ([]*affiliatedorganisation.AffiliatedOrganisation, error)
	Create(ctx context.Context, org affiliatedorganisation.AffiliatedOrganisation) (*affiliatedorganisation.AffiliatedOrganisation, error)
	Update(ctx context.Context, id uuid.UUID, org affiliatedorganisation.AffiliatedOrganisation) (*affiliatedorganisation.AffiliatedOrganisation, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]affiliatedorganisation.AffiliatedOrganisation, error)
}
