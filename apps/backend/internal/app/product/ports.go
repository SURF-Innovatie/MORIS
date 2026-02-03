package product

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/product"
	"github.com/google/uuid"
)

type Repository interface {
	Get(ctx context.Context, id uuid.UUID) (*product.Product, error)
	List(ctx context.Context) ([]*product.Product, error)
	ListByAuthorPersonID(ctx context.Context, personID uuid.UUID) ([]*product.Product, error)
	Create(ctx context.Context, p product.Product) (*product.Product, error)
	Update(ctx context.Context, id uuid.UUID, p product.Product) (*product.Product, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
