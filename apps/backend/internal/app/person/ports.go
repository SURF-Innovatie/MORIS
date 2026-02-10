package person

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/identity"
	"github.com/google/uuid"
)

type repository interface {
	Create(ctx context.Context, p identity.Person) (*identity.Person, error)
	Get(ctx context.Context, id uuid.UUID) (*identity.Person, error)
	Update(ctx context.Context, id uuid.UUID, p identity.Person) (*identity.Person, error)
	List(ctx context.Context) ([]*identity.Person, error)
	GetByEmail(ctx context.Context, email string) (*identity.Person, error)
	GetByORCID(ctx context.Context, orcid string) (*identity.Person, error)
	Search(ctx context.Context, query string, limit int) ([]identity.Person, error)
	SetORCID(ctx context.Context, personID uuid.UUID, orcidID string) error
	ClearORCID(ctx context.Context, personID uuid.UUID) error
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]identity.Person, error)
}
