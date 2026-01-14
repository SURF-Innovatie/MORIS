package person

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, p entities.Person) (*entities.Person, error)
	Get(ctx context.Context, id uuid.UUID) (*entities.Person, error)
	Update(ctx context.Context, id uuid.UUID, p entities.Person) (*entities.Person, error)
	List(ctx context.Context) ([]*entities.Person, error)
	GetByEmail(ctx context.Context, email string) (*entities.Person, error)
	Search(ctx context.Context, query string, limit int) ([]entities.Person, error)
	SetORCID(ctx context.Context, personID uuid.UUID, orcidID string) error
	ClearORCID(ctx context.Context, personID uuid.UUID) error
}
