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

	projects := make([]*ProjectDetails, 0, len(ids))
	for _, id := range ids {
		proj, err := s.loader.Load(ctx, id)
		if err != nil {
			continue
		}

		details, err := s.buildProjectDetails(ctx, proj)
		if err != nil {
			return nil, err
		}

		projects = append(projects, details)
	}

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

	var pending []events.Event
	for _, e := range evts {
		if e.GetStatus() == "pending" {
			pending = append(pending, e)
		}
	}

	return pending, nil
}

func (s *service) buildProjectDetails(ctx context.Context, proj *entities.Project) (*ProjectDetails, error) {
	if proj == nil {
		return nil, errors.New("project is nil")
	}

	// Fetch People
	personIDSet := map[uuid.UUID]struct{}{}
	roleIDSet := map[uuid.UUID]struct{}{}

	for _, m := range proj.Members {
		personIDSet[m.PersonID] = struct{}{}
		roleIDSet[m.ProjectRoleID] = struct{}{}
	}
	personIDs := make([]uuid.UUID, 0, len(personIDSet))
	for id := range personIDSet {
		personIDs = append(personIDs, id)
	}
	roleIDs := make([]uuid.UUID, 0, len(roleIDSet))
	for id := range roleIDSet {
		roleIDs = append(roleIDs, id)
	}

	peopleMap, err := s.repo.PeopleByIDs(ctx, personIDs)
	if err != nil {
		return nil, err
	}

	rolesMap, err := s.repo.ProjectRolesByIDs(ctx, roleIDs)
	if err != nil {
		return nil, err
	}

	members := make([]entities.ProjectMemberDetail, 0, len(proj.Members))
	for _, m := range proj.Members {
		p, okP := peopleMap[m.PersonID]
		r, okR := rolesMap[m.ProjectRoleID]
		if okP && okR {
			members = append(members, entities.ProjectMemberDetail{
				Person: p,
				Role:   r,
			})
		}
	}

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
	for _, evt := range evts {
		log.Entries = append(log.Entries, entities.ChangeLogEntry{
			Event: evt.String(),
			At:    evt.OccurredAt(),
		})
	}

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
