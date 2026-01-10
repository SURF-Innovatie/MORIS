package projectrole

import (
	"context"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/internal/app/organisation"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type Service interface {
	EnsureDefaults(ctx context.Context) error
	Create(ctx context.Context, key, name string, orgNodeID uuid.UUID) (*entities.ProjectRole, error)
	ListAvailableForNode(ctx context.Context, orgNodeID uuid.UUID) ([]entities.ProjectRole, error)
	Delete(ctx context.Context, id uuid.UUID, orgNodeID uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.ProjectRole, error)
	UpdateAllowedEventTypes(ctx context.Context, roleID uuid.UUID, eventTypes []string) (*entities.ProjectRole, error)
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

	// Get all registered event types for default roles
	allEventTypes := events.GetRegisteredEventTypes()

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
				// Create with all event types allowed
				_, err := s.repo.CreateWithEventTypes(ctx, d.key, d.name, root.ID, allEventTypes)
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

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*entities.ProjectRole, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) UpdateAllowedEventTypes(ctx context.Context, roleID uuid.UUID, eventTypes []string) (*entities.ProjectRole, error) {
	return s.repo.UpdateAllowedEventTypes(ctx, roleID, eventTypes)
}
