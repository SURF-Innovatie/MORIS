package queries

import (
	"context"
	"errors"
	"sort"

	appauth "github.com/SURF-Innovatie/MORIS/internal/app/auth"
	"github.com/SURF-Innovatie/MORIS/internal/app/event"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/load"
	"github.com/SURF-Innovatie/MORIS/internal/app/user"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

// ErrNotFound is returned when a project does not exist (no events).
var ErrNotFound = errors.New("project not found")

type Service interface {
	GetProject(ctx context.Context, id uuid.UUID) (*ProjectDetails, error)
	GetAllProjects(ctx context.Context) ([]*ProjectDetails, error)
	GetChangeLog(ctx context.Context, id uuid.UUID) ([]events.DetailedEvent, error)
	GetPendingEvents(ctx context.Context, projectID uuid.UUID) ([]events.DetailedEvent, error)
	GetProjectRoles(ctx context.Context) ([]entities.ProjectRole, error)
	ListAvailableRoles(ctx context.Context, projectID uuid.UUID) ([]entities.ProjectRole, error)
	GetEvents(ctx context.Context, id uuid.UUID) ([]events.Event, error)
	GetAllowedEventTypes(ctx context.Context, projectID uuid.UUID) ([]string, error)
	GetProjectRolesByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]entities.ProjectRole, error)
	GetProductsByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]entities.Product, error)
	GetPeopleByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]entities.Person, error)
}

type service struct {
	repo        ProjectReadRepository
	eventSvc    event.Service
	loader      *load.Loader
	currentUser appauth.CurrentUserProvider
	roleRepo    ProjectRoleRepository
	userSvc     user.Service
}

func NewService(
	eventSvc event.Service,
	loader *load.Loader,
	repo ProjectReadRepository,
	roleRepo ProjectRoleRepository,
	currentUser appauth.CurrentUserProvider,
	userSvc user.Service,
) Service {
	return &service{
		eventSvc:    eventSvc,
		loader:      loader,
		repo:        repo,
		roleRepo:    roleRepo,
		currentUser: currentUser,
		userSvc:     userSvc,
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

func (s *service) GetPendingEvents(ctx context.Context, projectID uuid.UUID) ([]events.DetailedEvent, error) {
	evts, _, err := s.eventSvc.Load(ctx, projectID)
	if err != nil {
		return nil, err
	}

	pendingEvts := lo.Filter(evts, func(e events.Event, _ int) bool {
		return e.GetStatus() == "pending"
	})

	// Hydration
	var personIDs []uuid.UUID
	var roleIDs []uuid.UUID
	var productIDs []uuid.UUID
	var creatorUserIDs []uuid.UUID

	for _, e := range pendingEvts {
		if r, ok := e.(events.HasRelatedIDs); ok {
			ids := r.RelatedIDs()
			if ids.PersonID != nil {
				personIDs = append(personIDs, *ids.PersonID)
			}
			if ids.ProjectRoleID != nil {
				roleIDs = append(roleIDs, *ids.ProjectRoleID)
			}
			if ids.ProductID != nil {
				productIDs = append(productIDs, *ids.ProductID)
			}
		}
		creatorUserIDs = append(creatorUserIDs, e.CreatedByID())
	}

	personMap := make(map[uuid.UUID]entities.Person)
	if len(personIDs) > 0 {
		pm, err := s.repo.PeopleByIDs(ctx, personIDs)
		if err == nil {
			personMap = pm
		}
	}

	roleMap := make(map[uuid.UUID]entities.ProjectRole)
	if len(roleIDs) > 0 {
		rm, err := s.repo.ProjectRolesByIDs(ctx, roleIDs)
		if err == nil {
			roleMap = rm
		}
	}

	productMap := make(map[uuid.UUID]entities.Product)
	if len(productIDs) > 0 {
		pm, err := s.repo.ProductsByIDs(ctx, productIDs)
		if err == nil {
			productMap = make(map[uuid.UUID]entities.Product)
			for _, p := range pm {
				productMap[p.Id] = p
			}
		}
	}

	creatorMap := make(map[uuid.UUID]entities.Person)
	if len(creatorUserIDs) > 0 {
		cm, err := s.userSvc.GetPeopleByUserIDs(ctx, lo.Uniq(creatorUserIDs))
		if err == nil {
			creatorMap = cm
		}
	}

	return lo.Map(pendingEvts, func(e events.Event, _ int) events.DetailedEvent {
		de := events.DetailedEvent{Event: e}

		if r, ok := e.(events.HasRelatedIDs); ok {
			ids := r.RelatedIDs()
			if ids.PersonID != nil {
				if p, ok := personMap[*ids.PersonID]; ok {
					de.Person = &p
				}
			}
			if ids.ProjectRoleID != nil {
				if r, ok := roleMap[*ids.ProjectRoleID]; ok {
					de.ProjectRole = &r
				}
			}
			if ids.ProductID != nil {
				if p, ok := productMap[*ids.ProductID]; ok {
					de.Product = &p
				}
			}
		}

		if p, ok := creatorMap[e.CreatedByID()]; ok {
			de.Creator = &p
		}

		return de
	}), nil
}

func (s *service) buildProjectDetails(ctx context.Context, proj *entities.Project) (*ProjectDetails, error) {
	if proj == nil {
		return nil, errors.New("project is nil")
	}

	personIDs := lo.Uniq(lo.Map(proj.Members, func(m entities.ProjectMember, _ int) uuid.UUID {
		return m.PersonID
	}))
	roleIDs := lo.Uniq(lo.Map(proj.Members, func(m entities.ProjectMember, _ int) uuid.UUID {
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

	members := lo.FilterMap(proj.Members, func(m entities.ProjectMember, _ int) (entities.ProjectMemberDetail, bool) {
		p, okP := peopleMap[m.PersonID]
		r, okR := rolesMap[m.ProjectRoleID]
		if okP && okR {
			return entities.ProjectMemberDetail{
				Person: p,
				Role:   r,
			}, true
		}
		return entities.ProjectMemberDetail{}, false
	})

	products, err := s.repo.ProductsByIDs(ctx, proj.ProductIDs)
	if err != nil {
		return nil, err
	}

	org, err := s.repo.OrganisationNodeByID(ctx, proj.OwningOrgNodeID)
	if err != nil {
		return nil, err
	}

	return &ProjectDetails{
		Project:       *proj,
		OwningOrgNode: org,
		Members:       members,
		Products:      products,
	}, nil
}

func (s *service) GetChangeLog(ctx context.Context, id uuid.UUID) ([]events.DetailedEvent, error) {
	evts, _, err := s.eventSvc.Load(ctx, id)
	if err != nil {
		return nil, err
	}
	if len(evts) == 0 {
		return nil, ErrNotFound
	}

	// Filter to only approved events for changelog
	approvedEvts := lo.Filter(evts, func(e events.Event, _ int) bool {
		return e.GetStatus() == "approved"
	})

	// Hydration - collect IDs
	var personIDs []uuid.UUID
	var roleIDs []uuid.UUID
	var productIDs []uuid.UUID
	var creatorUserIDs []uuid.UUID

	for _, e := range approvedEvts {
		if r, ok := e.(events.HasRelatedIDs); ok {
			ids := r.RelatedIDs()
			if ids.PersonID != nil {
				personIDs = append(personIDs, *ids.PersonID)
			}
			if ids.ProjectRoleID != nil {
				roleIDs = append(roleIDs, *ids.ProjectRoleID)
			}
			if ids.ProductID != nil {
				productIDs = append(productIDs, *ids.ProductID)
			}
		}
		creatorUserIDs = append(creatorUserIDs, e.CreatedByID())
	}

	personMap := make(map[uuid.UUID]entities.Person)
	if len(personIDs) > 0 {
		pm, err := s.repo.PeopleByIDs(ctx, personIDs)
		if err == nil {
			personMap = pm
		}
	}

	roleMap := make(map[uuid.UUID]entities.ProjectRole)
	if len(roleIDs) > 0 {
		rm, err := s.repo.ProjectRolesByIDs(ctx, roleIDs)
		if err == nil {
			roleMap = rm
		}
	}

	productMap := make(map[uuid.UUID]entities.Product)
	if len(productIDs) > 0 {
		pm, err := s.repo.ProductsByIDs(ctx, productIDs)
		if err == nil {
			productMap = make(map[uuid.UUID]entities.Product)
			for _, p := range pm {
				productMap[p.Id] = p
			}
		}
	}

	creatorMap := make(map[uuid.UUID]entities.Person)
	if len(creatorUserIDs) > 0 {
		cm, err := s.userSvc.GetPeopleByUserIDs(ctx, lo.Uniq(creatorUserIDs))
		if err == nil {
			creatorMap = cm
		}
	}

	detailedEvents := lo.Map(approvedEvts, func(e events.Event, _ int) events.DetailedEvent {
		de := events.DetailedEvent{Event: e}

		if r, ok := e.(events.HasRelatedIDs); ok {
			ids := r.RelatedIDs()
			if ids.PersonID != nil {
				if p, ok := personMap[*ids.PersonID]; ok {
					de.Person = &p
				}
			}
			if ids.ProjectRoleID != nil {
				if r, ok := roleMap[*ids.ProjectRoleID]; ok {
					de.ProjectRole = &r
				}
			}
			if ids.ProductID != nil {
				if p, ok := productMap[*ids.ProductID]; ok {
					de.Product = &p
				}
			}
		}

		if p, ok := creatorMap[e.CreatedByID()]; ok {
			de.Creator = &p
		}

		return de
	})

	// Sort by time descending (most recent first)
	sort.Slice(detailedEvents, func(i, j int) bool {
		return detailedEvents[i].Event.OccurredAt().After(detailedEvents[j].Event.OccurredAt())
	})

	return detailedEvents, nil
}

func (s *service) GetProjectRoles(ctx context.Context) ([]entities.ProjectRole, error) {
	return s.roleRepo.List(ctx)
}

func (s *service) ListAvailableRoles(ctx context.Context, projectID uuid.UUID) ([]entities.ProjectRole, error) {
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

func (s *service) GetEvents(ctx context.Context, id uuid.UUID) ([]events.Event, error) {
	evts, _, err := s.eventSvc.Load(ctx, id)
	return evts, err
}

func (s *service) GetAllowedEventTypes(ctx context.Context, projectID uuid.UUID) ([]string, error) {
	u, err := s.currentUser.Current(ctx)
	if err != nil {
		return nil, err
	}

	if u.IsSysAdmin {
		return events.GetRegisteredEventTypes(), nil
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

func (s *service) GetProjectRolesByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]entities.ProjectRole, error) {
	return s.repo.ProjectRolesByIDs(ctx, ids)
}

func (s *service) GetProductsByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]entities.Product, error) {
	products, err := s.repo.ProductsByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	return lo.SliceToMap(products, func(p entities.Product) (uuid.UUID, entities.Product) {
		return p.Id, p
	}), nil
}

func (s *service) GetPeopleByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]entities.Person, error) {
	return s.repo.PeopleByIDs(ctx, ids)
}
