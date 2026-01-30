package hierarchy

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	AncestorIDsInclusive(ctx context.Context, nodeID uuid.UUID) ([]uuid.UUID, error)
	AncestorIDs(ctx context.Context, nodeID uuid.UUID) ([]uuid.UUID, error)
	IsAncestor(ctx context.Context, ancestorID, descendantID uuid.UUID) (bool, error)
}

type service struct {
	repo repository
}

func NewService(repo repository) Service {
	return &service{repo: repo}
}

func (s *service) AncestorIDsInclusive(ctx context.Context, nodeID uuid.UUID) ([]uuid.UUID, error) {
	return s.repo.AncestorIDs(ctx, nodeID) // repo returns inclusive
}

func (s *service) AncestorIDs(ctx context.Context, nodeID uuid.UUID) ([]uuid.UUID, error) {
	ids, err := s.repo.AncestorIDs(ctx, nodeID)
	if err != nil {
		return nil, err
	}

	out := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		if id != nodeID {
			out = append(out, id)
		}
	}
	return out, nil
}

func (s *service) IsAncestor(ctx context.Context, ancestorID, descendantID uuid.UUID) (bool, error) {
	return s.repo.IsAncestor(ctx, ancestorID, descendantID)
}
