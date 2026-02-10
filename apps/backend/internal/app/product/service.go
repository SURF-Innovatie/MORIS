package product

import (
	"context"
	"fmt"
	"strings"

	"github.com/SURF-Innovatie/MORIS/internal/api/dto"
	"github.com/SURF-Innovatie/MORIS/internal/app/doi"
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
	GetOrCreateFromDOI(ctx context.Context, authorPersonID uuid.UUID, doi string) (*product.Product, bool, error)
}

type service struct {
	repo   Repository
	doiSvc doi.Service
}

func NewService(repo Repository, doiSvc doi.Service) Service {
	return &service{repo: repo, doiSvc: doiSvc}
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
	doiStr := strings.TrimSpace(w.DOI)
	if doiStr == "" {
		return nil, false, fmt.Errorf("work.doi is required")
	}

	existing, err := s.repo.GetByDOI(ctx, doiStr)
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
		DOI:            doiStr,
		AuthorPersonID: authorPersonID,
	}

	created, err := s.repo.Create(ctx, p)
	if err != nil {
		return nil, false, err
	}
	return created, true, nil
}

func (s *service) GetOrCreateFromDOI(ctx context.Context, authorPersonID uuid.UUID, doiStr string) (*product.Product, bool, error) {
	doiStr = strings.TrimSpace(doiStr)
	if doiStr == "" {
		return nil, false, fmt.Errorf("doi is required")
	}

	existing, err := s.repo.GetByDOI(ctx, doiStr)
	if err != nil {
		return nil, false, err
	}
	if existing != nil {
		return existing, false, nil
	}

	if s.doiSvc == nil {
		return nil, false, fmt.Errorf("doi service not configured")
	}

	work, err := s.doiSvc.Resolve(ctx, doiStr)
	if err != nil {
		return nil, false, err
	}

	return s.CreateOrGetFromWork(ctx, authorPersonID, work)
}
