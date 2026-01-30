package person

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, p entities.Person) (*entities.Person, error)
	Get(ctx context.Context, id uuid.UUID) (*entities.Person, error)
	Update(ctx context.Context, id uuid.UUID, p entities.Person) (*entities.Person, error)
	List(ctx context.Context) ([]*entities.Person, error)
	GetByEmail(ctx context.Context, email string) (*entities.Person, error)
	Search(ctx context.Context, query string, limit int) ([]entities.Person, error)
	SetORCID(ctx context.Context, personID uuid.UUID, orcidID string) error
	ClearORCID(ctx context.Context, personID uuid.UUID) error
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]entities.Person, error)
}

type service struct {
	repo repository
}

func NewService(repo repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, p entities.Person) (*entities.Person, error) {
	return s.repo.Create(ctx, p)
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*entities.Person, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) Update(ctx context.Context, id uuid.UUID, p entities.Person) (*entities.Person, error) {
	return s.repo.Update(ctx, id, p)
}

func (s *service) List(ctx context.Context) ([]*entities.Person, error) {
	return s.repo.List(ctx)
}

func (s *service) GetByEmail(ctx context.Context, email string) (*entities.Person, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *service) Search(ctx context.Context, query string, limit int) ([]entities.Person, error) {
	return s.repo.Search(ctx, query, limit)
}

func (s *service) SetORCID(ctx context.Context, personID uuid.UUID, orcidID string) error {
	return s.repo.SetORCID(ctx, personID, orcidID)
}

func (s *service) ClearORCID(ctx context.Context, personID uuid.UUID) error {
	return s.repo.ClearORCID(ctx, personID)
}

func (s *service) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]entities.Person, error) {
	return s.repo.GetByIDs(ctx, ids)
}
