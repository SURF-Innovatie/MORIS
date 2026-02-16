package queries

import (
	"context"
	"errors"
	"sort"

	appauth "github.com/SURF-Innovatie/MORIS/internal/app/auth"
	"github.com/SURF-Innovatie/MORIS/internal/app/event"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/load"
	"github.com/SURF-Innovatie/MORIS/internal/app/user"
	"github.com/SURF-Innovatie/MORIS/internal/domain/affiliatedorganisation"
	"github.com/SURF-Innovatie/MORIS/internal/domain/identity"
	"github.com/SURF-Innovatie/MORIS/internal/domain/product"
	"github.com/SURF-Innovatie/MORIS/internal/domain/project"
	events2 "github.com/SURF-Innovatie/MORIS/internal/domain/project/events"
	"github.com/SURF-Innovatie/MORIS/internal/domain/project/events/hydrator"
	"github.com/SURF-Innovatie/MORIS/internal/domain/project/role"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

// ErrNotFound is returned when a project does not exist (no events).
var ErrNotFound = errors.New("project not found")

type Service interface {
	GetProject(ctx context.Context, id uuid.UUID) (*ProjectDetails, error)
	GetAllProjects(ctx context.Context) ([]*ProjectDetails, error)
	GetChangeLog(ctx context.Context, id uuid.UUID) ([]events2.DetailedEvent, error)
	GetPendingEvents(ctx context.Context, projectID uuid.UUID) ([]events2.DetailedEvent, error)
	GetProjectRoles(ctx context.Context) ([]role.ProjectRole, error)
	ListAvailableRoles(ctx context.Context, projectID uuid.UUID) ([]role.ProjectRole, error)
	GetEvents(ctx context.Context, id uuid.UUID) ([]events2.Event, error)
	GetAllowedEventTypes(ctx context.Context, projectID uuid.UUID) ([]string, error)
	GetProjectRolesByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]role.ProjectRole, error)
	GetProductsByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]product.Product, error)
	GetPeopleByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]identity.Person, error)
}

type service struct {
	repo        ProjectReadRepository
	eventSvc    event.Service
	loader      *load.Loader
	currentUser appauth.CurrentUserProvider
	roleRepo    ProjectRoleRepository
	userSvc     user.Service
	hydrator    *hydrator.Hydrator
}

func NewService(
	eventSvc event.Service,
	loader *load.Loader,
	repo ProjectReadRepository,
	roleRepo ProjectRoleRepository,
	currentUser appauth.CurrentUserProvider,
	userSvc user.Service,
	h *hydrator.Hydrator,
) Service {
	return &service{
		eventSvc:    eventSvc,
		loader:      loader,
		repo:        repo,
		roleRepo:    roleRepo,
		currentUser: currentUser,
		userSvc:     userSvc,
		hydrator:    h,
	}
}

func (s *service) GetProject(ctx context.Context, id uuid.UUID) (*ProjectDetails, error) {
	proj, err := s.loader.Load(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.buildProjectDetails(ctx, proj)
}

func (s *service) GetAllProjects(ctx context.Context) ([]*ProjectDetails, error) {
	u, err := s.currentUser.Current(ctx)
	if err != nil {
		return nil, err
	}

	var ids []uuid.UUID

	// Sysadmins can see all projects
	if u.IsSysAdmin {
		ids, err = s.repo.ProjectIDsStarted(ctx)
	} else {
		ids, err = s.repo.ProjectIDsForPerson(ctx, u.PersonID)
	}
	if err != nil {
		return nil, err
	}

	projects := lo.FilterMap(ids, func(id uuid.UUID, _ int) (*ProjectDetails, bool) {
		proj, err := s.loader.Load(ctx, id)
		if err != nil {
			return nil, false
		}

		details, err := s.buildProjectDetails(ctx, proj)
		if err != nil {
			return nil, false
		}

		return details, true
	})

	// Sort by title for consistency
	sort.Slice(projects, func(i, j int) bool {
		return projects[i].Project.Title < projects[j].Project.Title
	})

	return projects, nil
}

func (s *service) GetPendingEvents(ctx context.Context, projectID uuid.UUID) ([]events2.DetailedEvent, error) {
	evts, _, err := s.eventSvc.Load(ctx, projectID)
	if err != nil {
		return nil, err
	}

	pendingEvts := lo.Filter(evts, func(e events2.Event, _ int) bool {
		return e.GetStatus() == "pending"
	})

	return s.hydrator.HydrateMany(ctx, pendingEvts), nil
}

func (s *service) buildProjectDetails(ctx context.Context, proj *project.Project) (*ProjectDetails, error) {
	if proj == nil {
		return nil, errors.New("project is nil")
	}

	personIDs := lo.Uniq(lo.Map(proj.Members, func(m project.Member, _ int) uuid.UUID {
		return m.PersonID
	}))
	roleIDs := lo.Uniq(lo.Map(proj.Members, func(m project.Member, _ int) uuid.UUID {
		return m.ProjectRoleID
	}))

	peopleMap, err := s.repo.PeopleByIDs(ctx, personIDs)
	if err != nil {
		return nil, err
	}

	rolesMap, err := s.repo.ProjectRolesByIDs(ctx, roleIDs)
	if err != nil {
		return nil, err
	}

	members := lo.FilterMap(proj.Members, func(m project.Member, _ int) (project.MemberDetail, bool) {
		p, okP := peopleMap[m.PersonID]
		r, okR := rolesMap[m.ProjectRoleID]
		if okP && okR {
			return project.MemberDetail{
				Person: p,
				Role:   r,
			}, true
		}
		return project.MemberDetail{}, false
	})

	products, err := s.repo.ProductsByIDs(ctx, proj.ProductIDs)
	if err != nil {
		return nil, err
	}

	org, err := s.repo.OrganisationNodeByID(ctx, proj.OwningOrgNodeID)
	if err != nil {
		return nil, err
	}

	// Load affiliated organisations
	affiliatedOrgsMap, err := s.repo.GetAffiliatedOrganisationsByIDs(ctx, proj.AffiliatedOrganisationIDs)
	if err != nil {
		return nil, err
	}
	affiliatedOrgs := make([]affiliatedorganisation.AffiliatedOrganisation, 0, len(affiliatedOrgsMap))
	for _, org := range affiliatedOrgsMap {
		affiliatedOrgs = append(affiliatedOrgs, org)
	}

	return &ProjectDetails{
		Project:                 *proj,
		OwningOrgNode:           org,
		Members:                 members,
		Products:                products,
		AffiliatedOrganisations: affiliatedOrgs,
	}, nil
}

func (s *service) GetChangeLog(ctx context.Context, id uuid.UUID) ([]events2.DetailedEvent, error) {
	evts, _, err := s.eventSvc.Load(ctx, id)
	if err != nil {
		return nil, err
	}
	if len(evts) == 0 {
		return nil, ErrNotFound
	}

	// Filter to only approved events for changelog
	approvedEvts := lo.Filter(evts, func(e events2.Event, _ int) bool {
		return e.GetStatus() == "approved"
	})

	detailedEvents := s.hydrator.HydrateMany(ctx, approvedEvts)

	// Sort by time descending (most recent first)
	sort.Slice(detailedEvents, func(i, j int) bool {
		return detailedEvents[i].Event.OccurredAt().After(detailedEvents[j].Event.OccurredAt())
	})

	return detailedEvents, nil
}

func (s *service) GetProjectRoles(ctx context.Context) ([]role.ProjectRole, error) {
	return s.roleRepo.List(ctx)
}

func (s *service) ListAvailableRoles(ctx context.Context, projectID uuid.UUID) ([]role.ProjectRole, error) {
	proj, err := s.loader.Load(ctx, projectID)
	if err != nil {
		return nil, err
	}

	ancestors, err := s.repo.ListAncestors(ctx, proj.OwningOrgNodeID)
	if err != nil {
		return nil, err
	}

	return s.roleRepo.ListByOrgIDs(ctx, ancestors)
}

func (s *service) GetEvents(ctx context.Context, id uuid.UUID) ([]events2.Event, error) {
	evts, _, err := s.eventSvc.Load(ctx, id)
	return evts, err
}

func (s *service) GetAllowedEventTypes(ctx context.Context, projectID uuid.UUID) ([]string, error) {
	u, err := s.currentUser.Current(ctx)
	if err != nil {
		return nil, err
	}

	if u.IsSysAdmin {
		return events2.GetRegisteredEventTypes(), nil
	}

	// Load the project to find user's role
	proj, err := s.loader.Load(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// Find user's role in this project
	var userRoleID *uuid.UUID
	for _, m := range proj.Members {
		if m.PersonID == u.PersonID {
			userRoleID = &m.ProjectRoleID
			break
		}
	}

	// If user is not a member, return empty list
	if userRoleID == nil {
		return []string{}, nil
	}

	// Get the role details
	rolesMap, err := s.repo.ProjectRolesByIDs(ctx, []uuid.UUID{*userRoleID})
	if err != nil {
		return nil, err
	}

	role, ok := rolesMap[*userRoleID]
	if !ok {
		return []string{}, nil
	}

	// Return allowed event types (or empty if none configured)
	return role.AllowedEventTypes, nil
}

func (s *service) GetProjectRolesByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]role.ProjectRole, error) {
	return s.repo.ProjectRolesByIDs(ctx, ids)
}

func (s *service) GetProductsByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]product.Product, error) {
	products, err := s.repo.ProductsByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	return lo.SliceToMap(products, func(p product.Product) (uuid.UUID, product.Product) {
		return p.Id, p
	}), nil
}

func (s *service) GetPeopleByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]identity.Person, error) {
	return s.repo.PeopleByIDs(ctx, ids)
}
