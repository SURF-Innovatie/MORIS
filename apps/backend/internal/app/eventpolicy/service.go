package eventpolicy

import (
	"context"

	organisationhierarchy "github.com/SURF-Innovatie/MORIS/internal/app/organisation/hierarchy"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

// Service provides event policy management and evaluation
type Service interface {
	Create(ctx context.Context, policy entities.EventPolicy) (*entities.EventPolicy, error)
	Update(ctx context.Context, id uuid.UUID, policy entities.EventPolicy) (*entities.EventPolicy, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.EventPolicy, error)

	ListForOrgNode(ctx context.Context, orgNodeID uuid.UUID, includeInherited bool) ([]entities.EventPolicy, error)
	ListForProject(ctx context.Context, projectID uuid.UUID, owningOrgNodeID uuid.UUID, includeInherited bool) ([]entities.EventPolicy, error)
}

type service struct {
	repo            repository
	orgHierarchySvc organisationhierarchy.Service
}

// NewService creates a new event policy service
func NewService(repo repository, orgRbacSvc organisationhierarchy.Service) Service {
	return &service{
		repo:            repo,
		orgHierarchySvc: orgRbacSvc,
	}
}

func (s *service) Create(ctx context.Context, policy entities.EventPolicy) (*entities.EventPolicy, error) {
	return s.repo.Create(ctx, policy)
}

func (s *service) Update(ctx context.Context, id uuid.UUID, policy entities.EventPolicy) (*entities.EventPolicy, error) {
	return s.repo.Update(ctx, id, policy)
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*entities.EventPolicy, error) {
	return s.repo.GetByID(ctx, id)
}

// ListForOrgNode lists policies for an org node, optionally including inherited ones
func (s *service) ListForOrgNode(ctx context.Context, orgNodeID uuid.UUID, includeInherited bool) ([]entities.EventPolicy, error) {
	var ancestorIDs []uuid.UUID
	if includeInherited {
		var err error
		ancestorIDs, err = s.orgHierarchySvc.AncestorIDs(ctx, orgNodeID)
		if err != nil {
			return nil, err
		}
	}

	policies, err := s.repo.ListForOrgNode(ctx, orgNodeID, ancestorIDs)
	if err != nil {
		return nil, err
	}

	// Mark inherited policies and set source info
	for i := range policies {
		if policies[i].OrgNodeID != nil && *policies[i].OrgNodeID != orgNodeID {
			policies[i].Inherited = true
			policies[i].SourceOrgNodeID = policies[i].OrgNodeID
		}
	}

	return policies, nil
}

// ListForProject lists policies for a project, optionally including inherited org policies
func (s *service) ListForProject(ctx context.Context, projectID uuid.UUID, owningOrgNodeID uuid.UUID, includeInherited bool) ([]entities.EventPolicy, error) {
	// Get project's own policies
	projectPolicies, err := s.repo.ListForProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if !includeInherited {
		return projectPolicies, nil
	}

	// Get inherited org node policies (including all ancestors)
	ancestorIDs, err := s.orgHierarchySvc.AncestorIDs(ctx, owningOrgNodeID)
	if err != nil {
		return nil, err
	}

	// Include the owning org node itself
	allOrgIDs := append([]uuid.UUID{owningOrgNodeID}, ancestorIDs...)

	orgPolicies, err := s.repo.ListForOrgNode(ctx, owningOrgNodeID, allOrgIDs)
	if err != nil {
		return nil, err
	}

	// Mark all org policies as inherited
	for i := range orgPolicies {
		orgPolicies[i].Inherited = true
		orgPolicies[i].SourceOrgNodeID = orgPolicies[i].OrgNodeID
	}

	// Combine: project policies first, then org policies
	return append(projectPolicies, orgPolicies...), nil
}
