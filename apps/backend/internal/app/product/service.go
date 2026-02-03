package product

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/product"
	"github.com/google/uuid"
)

type Service interface {
	Get(ctx context.Context, id uuid.UUID) (*product.Product, error)
	GetAll(ctx context.Context) ([]*product.Product, error)
	GetAllForUser(ctx context.Context, personID uuid.UUID) ([]*product.Product, error)
	Create(ctx context.Context, p product.Product) (*product.Product, error)
	Update(ctx context.Context, id uuid.UUID, p product.Product) (*product.Product, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*product.Product, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) GetAll(ctx context.Context) ([]*product.Product, error) {
	return s.repo.List(ctx)
}

func (s *service) GetAllForUser(ctx context.Context, personID uuid.UUID) ([]*product.Product, error) {
	return s.repo.ListByAuthorPersonID(ctx, personID)
}

func (s *service) Create(ctx context.Context, p product.Product) (*product.Product, error) {
	return s.repo.Create(ctx, p)
}

func (s *service) Update(ctx context.Context, id uuid.UUID, p product.Product) (*product.Product, error) {
	return s.repo.Update(ctx, id, p)
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
