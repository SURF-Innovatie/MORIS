package product

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type Repository interface {
	Get(ctx context.Context, id uuid.UUID) (*entities.Product, error)
	List(ctx context.Context) ([]*entities.Product, error)
	ListByAuthorPersonID(ctx context.Context, personID uuid.UUID) ([]*entities.Product, error)
	Create(ctx context.Context, p entities.Product) (*entities.Product, error)
	Update(ctx context.Context, id uuid.UUID, p entities.Product) (*entities.Product, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
