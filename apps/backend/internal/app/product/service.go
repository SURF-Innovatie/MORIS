package product

import (
	"context"
	"fmt"
	"strings"

	"github.com/SURF-Innovatie/MORIS/internal/api/dto"
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
	GetByDOI(ctx context.Context, doi string) (*product.Product, error)
	CreateOrGetFromWork(ctx context.Context, authorPersonID uuid.UUID, work *dto.Work) (*product.Product, bool, error)
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

func (s *service) GetByDOI(ctx context.Context, doi string) (*product.Product, error) {
	doi = strings.TrimSpace(doi)
	if doi == "" {
		return nil, nil
	}
	return s.repo.GetByDOI(ctx, doi)
}

func (s *service) CreateOrGetFromWork(ctx context.Context, authorPersonID uuid.UUID, w *dto.Work) (*product.Product, bool, error) {
	if w == nil {
		return nil, false, fmt.Errorf("work is required")
	}
	doi := strings.TrimSpace(w.DOI)
	if doi == "" {
		return nil, false, fmt.Errorf("work.doi is required")
	}

	existing, err := s.repo.GetByDOI(ctx, doi)
	if err != nil {
		return nil, false, err
	}
	if existing != nil {
		return existing, false, nil
	}

	p := product.Product{
		Name:           w.Title,
		Type:           w.Type,
		Language:       "en", // TODO: derive from metadata if available
		DOI:            doi,
		AuthorPersonID: authorPersonID,
	}
	created, err := s.repo.Create(ctx, p)
	if err != nil {
		return nil, false, err
	}
	return created, true, nil
}
