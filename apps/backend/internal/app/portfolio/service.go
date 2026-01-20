package portfolio

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type Service interface {
	GetForPerson(ctx context.Context, personID uuid.UUID) (*entities.Portfolio, error)
	UpdateForPerson(ctx context.Context, personID uuid.UUID, in entities.Portfolio) (*entities.Portfolio, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetForPerson(ctx context.Context, personID uuid.UUID) (*entities.Portfolio, error) {
	return s.repo.GetByPersonID(ctx, personID)
}

func (s *service) UpdateForPerson(ctx context.Context, personID uuid.UUID, in entities.Portfolio) (*entities.Portfolio, error) {
	in.PersonID = personID
	return s.repo.Upsert(ctx, in)
}
