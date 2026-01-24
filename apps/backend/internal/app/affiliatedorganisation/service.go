package affiliatedorganisation

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/affiliatedorganisation"
	"github.com/google/uuid"
)

// Service defines the business logic interface for AffiliatedOrganisation.
type Service interface {
	Get(ctx context.Context, id uuid.UUID) (*affiliatedorganisation.AffiliatedOrganisation, error)
	GetAll(ctx context.Context) ([]*affiliatedorganisation.AffiliatedOrganisation, error)
	Create(ctx context.Context, org affiliatedorganisation.AffiliatedOrganisation) (*affiliatedorganisation.AffiliatedOrganisation, error)
	Update(ctx context.Context, id uuid.UUID, org affiliatedorganisation.AffiliatedOrganisation) (*affiliatedorganisation.AffiliatedOrganisation, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]affiliatedorganisation.AffiliatedOrganisation, error)
}

type service struct {
	repo Repository
}

// NewService creates a new AffiliatedOrganisation service.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*affiliatedorganisation.AffiliatedOrganisation, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) GetAll(ctx context.Context) ([]*affiliatedorganisation.AffiliatedOrganisation, error) {
	return s.repo.List(ctx)
}

func (s *service) Create(ctx context.Context, org affiliatedorganisation.AffiliatedOrganisation) (*affiliatedorganisation.AffiliatedOrganisation, error) {
	return s.repo.Create(ctx, org)
}

func (s *service) Update(ctx context.Context, id uuid.UUID, org affiliatedorganisation.AffiliatedOrganisation) (*affiliatedorganisation.AffiliatedOrganisation, error) {
	return s.repo.Update(ctx, id, org)
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) GetByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]affiliatedorganisation.AffiliatedOrganisation, error) {
	return s.repo.GetByIDs(ctx, ids)
}
