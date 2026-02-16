package portfolio

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/portfolio"
	"github.com/google/uuid"
)

type Service interface {
	GetForPerson(ctx context.Context, personID uuid.UUID) (*portfolio.Portfolio, error)
	UpdateForPerson(ctx context.Context, personID uuid.UUID, in portfolio.Portfolio) (*portfolio.Portfolio, error)
	TrackProjectAccess(ctx context.Context, personID uuid.UUID, projectID uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetForPerson(ctx context.Context, personID uuid.UUID) (*portfolio.Portfolio, error) {
	return s.repo.GetByPersonID(ctx, personID)
}

func (s *service) UpdateForPerson(ctx context.Context, personID uuid.UUID, in portfolio.Portfolio) (*portfolio.Portfolio, error) {
	in.PersonID = personID
	return s.repo.Upsert(ctx, in)
}

func (s *service) TrackProjectAccess(ctx context.Context, personID uuid.UUID, projectID uuid.UUID) error {
	// Get existing portfolio
	p, err := s.repo.GetByPersonID(ctx, personID)
	if err != nil {
		return err
	}

	// Initialize if nil
	if p.RecentProjectIDs == nil {
		p.RecentProjectIDs = []uuid.UUID{}
	}

	// Remove projectID if it already exists (to move it to front)
	filtered := make([]uuid.UUID, 0, len(p.RecentProjectIDs))
	for _, id := range p.RecentProjectIDs {
		if id != projectID {
			filtered = append(filtered, id)
		}
	}

	// Add projectID to front
	p.RecentProjectIDs = append([]uuid.UUID{projectID}, filtered...)

	// Keep only last 10 recent projects
	if len(p.RecentProjectIDs) > 10 {
		p.RecentProjectIDs = p.RecentProjectIDs[:10]
	}

	// Save
	_, err = s.repo.Upsert(ctx, *p)
	return err
}
