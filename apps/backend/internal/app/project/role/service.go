package role

import (
	"context"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/internal/app/organisation"
	organisationhierarchy "github.com/SURF-Innovatie/MORIS/internal/app/organisation/hierarchy"
	"github.com/SURF-Innovatie/MORIS/internal/domain/project/events"
	"github.com/SURF-Innovatie/MORIS/internal/domain/project/role"
	"github.com/google/uuid"
)

type Service interface {
	EnsureDefaults(ctx context.Context) error
	Create(ctx context.Context, key, name string, orgNodeID uuid.UUID) (*role.ProjectRole, error)
	ListAvailableForNode(ctx context.Context, orgNodeID uuid.UUID) ([]role.ProjectRole, error)
	Delete(ctx context.Context, id uuid.UUID, orgNodeID uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*role.ProjectRole, error)
	UpdateAllowedEventTypes(ctx context.Context, roleID uuid.UUID, eventTypes []string) (*role.ProjectRole, error)
}

type service struct {
	repo            Repository
	orgSvc          organisation.Service
	orgHierarchySvc organisationhierarchy.Service
}

func NewService(repo Repository, orgSvc organisation.Service, orgHierarchySvc organisationhierarchy.Service) Service {
	return &service{repo: repo, orgSvc: orgSvc, orgHierarchySvc: orgHierarchySvc}
}

func (s *service) EnsureDefaults(ctx context.Context) error {
	roots, err := s.orgSvc.ListRoots(ctx)
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

func (s *service) Create(ctx context.Context, key, name string, orgNodeID uuid.UUID) (*role.ProjectRole, error) {
	return s.repo.CreateOrRestore(ctx, key, name, orgNodeID)
}

func (s *service) Delete(ctx context.Context, id uuid.UUID, orgNodeID uuid.UUID) error {
	return s.repo.Delete(ctx, id, orgNodeID)
}

func (s *service) ListAvailableForNode(ctx context.Context, orgNodeID uuid.UUID) ([]role.ProjectRole, error) {
	ids, err := s.orgHierarchySvc.AncestorIDs(ctx, orgNodeID)
	if err != nil {
		return nil, err
	}

	return s.repo.ListByOrgIDs(ctx, ids)
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*role.ProjectRole, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) UpdateAllowedEventTypes(ctx context.Context, roleID uuid.UUID, eventTypes []string) (*role.ProjectRole, error) {
	return s.repo.UpdateAllowedEventTypes(ctx, roleID, eventTypes)
}
