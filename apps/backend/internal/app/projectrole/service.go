package projectrole

import (
	"context"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/internal/app/organisation"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type Service interface {
	EnsureDefaults(ctx context.Context) error
	Create(ctx context.Context, key, name string, orgNodeID uuid.UUID) (*entities.ProjectRole, error)
	ListAvailableForNode(ctx context.Context, orgNodeID uuid.UUID) ([]entities.ProjectRole, error)
	Delete(ctx context.Context, id uuid.UUID, orgNodeID uuid.UUID) error
}

type service struct {
	repo    Repository
	orgRepo organisation.Repository
}

func NewService(repo Repository, orgRepo organisation.Repository) Service {
	return &service{repo: repo, orgRepo: orgRepo}
}

func (s *service) EnsureDefaults(ctx context.Context) error {
	roots, err := s.orgRepo.ListRoots(ctx)
	if err != nil {
		return fmt.Errorf("listing root nodes: %w", err)
	}

	defs := []struct {
		key  string
		name string
	}{
		{key: "contributor", name: "Contributor"},
		{key: "lead", name: "Project lead"},
	}

	for _, root := range roots {
		for _, d := range defs {
			exists, err := s.repo.Exists(ctx, d.key, root.ID)
			if err != nil {
				return fmt.Errorf("checking existence of role %s for root %s: %w", d.key, root.ID, err)
			}

			if !exists {
				_, err := s.repo.Create(ctx, d.key, d.name, root.ID)
				if err != nil {
					return fmt.Errorf("creating default role %s for root %s: %w", d.key, root.ID, err)
				}
			} else {
				err := s.repo.Unarchive(ctx, d.key, root.ID)
				if err != nil {
					return fmt.Errorf("ensuring default role %s is unarchived for root %s: %w", d.key, root.ID, err)
				}
			}
		}
	}

	return nil
}

func (s *service) Create(ctx context.Context, key, name string, orgNodeID uuid.UUID) (*entities.ProjectRole, error) {
	return s.repo.CreateOrRestore(ctx, key, name, orgNodeID)
}

func (s *service) Delete(ctx context.Context, id uuid.UUID, orgNodeID uuid.UUID) error {
	return s.repo.Delete(ctx, id, orgNodeID)
}

func (s *service) ListAvailableForNode(ctx context.Context, orgNodeID uuid.UUID) ([]entities.ProjectRole, error) {
	closures, err := s.orgRepo.ListClosuresByDescendant(ctx, orgNodeID)
	if err != nil {
		return nil, err
	}

	ids := lo.Map(closures, func(c entities.OrganisationNodeClosure, _ int) uuid.UUID {
		return c.AncestorID
	})

	return s.repo.ListByOrgIDs(ctx, ids)
}
