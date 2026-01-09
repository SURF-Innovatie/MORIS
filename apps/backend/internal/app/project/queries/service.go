package queries

import (
	"context"
	"errors"
	"sort"

	appauth "github.com/SURF-Innovatie/MORIS/internal/app/auth"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/load"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/eventstore"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

// ErrNotFound is returned when a project does not exist (no events).
var ErrNotFound = errors.New("project not found")

type Service interface {
	GetProject(ctx context.Context, id uuid.UUID) (*ProjectDetails, error)
	GetAllProjects(ctx context.Context) ([]*ProjectDetails, error)
	GetChangeLog(ctx context.Context, id uuid.UUID) (*entities.ChangeLog, error)
	GetPendingEvents(ctx context.Context, projectID uuid.UUID) ([]events.Event, error)
	GetProjectRoles(ctx context.Context) ([]entities.ProjectRole, error)
	ListAvailableRoles(ctx context.Context, projectID uuid.UUID) ([]entities.ProjectRole, error)
	GetEvents(ctx context.Context, id uuid.UUID) ([]events.Event, error)
}

type service struct {
	repo        ProjectReadRepository
	es          eventstore.Store
	loader      *load.Loader
	currentUser appauth.CurrentUserProvider
	roleRepo    ProjectRoleRepository
}

func NewService(
	es eventstore.Store,
	loader *load.Loader,
	repo ProjectReadRepository,
	roleRepo ProjectRoleRepository,
	currentUser appauth.CurrentUserProvider,
) Service {
	return &service{
		es:          es,
		loader:      loader,
		repo:        repo,
		roleRepo:    roleRepo,
		currentUser: currentUser,
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

	ids, err := s.repo.ProjectIDsForPerson(ctx, u.PersonID())
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

func (s *service) GetPendingEvents(ctx context.Context, projectID uuid.UUID) ([]events.Event, error) {
	evts, _, err := s.es.Load(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return lo.Filter(evts, func(e events.Event, _ int) bool {
		return e.GetStatus() == "pending"
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

func (s *service) GetChangeLog(ctx context.Context, id uuid.UUID) (*entities.ChangeLog, error) {
	evts, _, err := s.es.Load(ctx, id)
	if err != nil {
		return nil, err
	}
	if len(evts) == 0 {
		return nil, ErrNotFound
	}

	var log entities.ChangeLog
	log.Entries = lo.Map(evts, func(evt events.Event, _ int) entities.ChangeLogEntry {
		return entities.ChangeLogEntry{
			Event: evt.String(),
			At:    evt.OccurredAt(),
		}
	})

	sort.Slice(log.Entries, func(i, j int) bool {
		return log.Entries[i].At.After(log.Entries[j].At)
	})

	return &log, nil
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
	evts, _, err := s.es.Load(ctx, id)
	return evts, err
}
