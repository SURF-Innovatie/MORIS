package project

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/SURF-Innovatie/MORIS/ent"
	en "github.com/SURF-Innovatie/MORIS/ent/event"
	organisationent "github.com/SURF-Innovatie/MORIS/ent/organisationnode"
	personent "github.com/SURF-Innovatie/MORIS/ent/person"
	productent "github.com/SURF-Innovatie/MORIS/ent/product"
	entprojectrole "github.com/SURF-Innovatie/MORIS/ent/projectrole"
	"github.com/SURF-Innovatie/MORIS/internal/app/commandbus"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/SURF-Innovatie/MORIS/internal/domain/projection"
	"github.com/SURF-Innovatie/MORIS/internal/event"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/eventstore"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

// ErrNotFound is returned when a project does not exist (no events).
var ErrNotFound = errors.New("project not found")

type Service interface {
	GetProject(ctx context.Context, id uuid.UUID) (*entities.ProjectDetails, error)
	GetAllProjects(ctx context.Context) ([]*entities.ProjectDetails, error)
	StartProject(ctx context.Context, params StartProjectParams) (*entities.ProjectDetails, error)
	UpdateProject(ctx context.Context, id uuid.UUID, params UpdateProjectParams) (*entities.ProjectDetails, error)
	AddPerson(ctx context.Context, projectID uuid.UUID, personID uuid.UUID) (*entities.ProjectDetails, error)
	RemovePerson(ctx context.Context, projectID uuid.UUID, personID uuid.UUID) (*entities.ProjectDetails, error)
	AddProduct(ctx context.Context, projectID uuid.UUID, productID uuid.UUID) (*entities.ProjectDetails, error)
	RemoveProduct(ctx context.Context, projectID uuid.UUID, productID uuid.UUID) (*entities.ProjectDetails, error)
	UpdateMemberRole(ctx context.Context, projectID uuid.UUID, personID uuid.UUID, roleKey string) (*entities.ProjectDetails, error)
	GetChangeLog(ctx context.Context, id uuid.UUID) (*entities.ChangeLog, error)
	GetPendingEvents(ctx context.Context, projectID uuid.UUID) ([]events.Event, error)
	GetProjectRoles(ctx context.Context) ([]entities.ProjectRole, error)
	WarmupCache(ctx context.Context) error
}

type service struct {
	cli    *ent.Client
	es     eventstore.Store
	evtSvc event.Service
	redis  *redis.Client

	exec *commandbus.Executor[entities.Project]
}

type StartProjectParams struct {
	Title           string
	Description     string
	OwningOrgNodeID uuid.UUID
	StartDate       time.Time
	EndDate         time.Time
}

type UpdateProjectParams struct {
	Title           string
	Description     string
	OwningOrgNodeID uuid.UUID
	StartDate       time.Time
	EndDate         time.Time
}

func NewService(es eventstore.Store, cli *ent.Client, evtSvc event.Service, rdb *redis.Client) Service {
	s := &service{es: es, cli: cli, evtSvc: evtSvc, redis: rdb}

	s.exec = commandbus.NewExecutor[entities.Project](
		es,
		evtSvc,
		Reducer{},
		NewReducer{},
	)

	evtSvc.RegisterStatusChangeHandler(s.onStatusChange)
	return s
}

func (s *service) GetProject(ctx context.Context, id uuid.UUID) (*entities.ProjectDetails, error) {
	proj, err := s.fromDb(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.buildProjectDetails(ctx, proj)
}

func (s *service) StartProject(ctx context.Context, params StartProjectParams) (*entities.ProjectDetails, error) {
	if params.Title == "" {
		return nil, errors.New("title is required")
	}

	projectID := uuid.New()

	user, err := currentUser(ctx, s.cli)
	if err != nil {
		return nil, err
	}

	adminRoleID, err := s.projectRoleID(ctx, "admin")
	if err != nil {
		return nil, err
	}

	proj, err := s.exec.Execute(ctx, projectID, func(ctx context.Context, cur *entities.Project) ([]events.Event, error) {
		var out []events.Event

		startEvt, err := events.DecideProjectStarted(
			projectID,
			user.ID,
			events.ProjectStartedInput{
				Title:           params.Title,
				Description:     params.Description,
				StartDate:       params.StartDate,
				EndDate:         params.EndDate,
				Members:         nil,
				OwningOrgNodeID: params.OwningOrgNodeID,
			},
			events.StatusApproved,
		)
		if err != nil {
			return nil, err
		}
		emit(cur, &out, startEvt)

		assignEvt, err := events.DecideProjectRoleAssigned(
			projectID,
			user.ID,
			cur,
			events.ProjectRoleAssignedInput{
				PersonID:      user.PersonID,
				ProjectRoleID: adminRoleID,
			},
			events.StatusApproved,
		)
		if err != nil {
			return nil, err
		}
		emit(cur, &out, assignEvt)

		return out, nil
	})

	if err != nil {
		return nil, err
	}

	_ = s.cacheProject(ctx, proj)

	return s.buildProjectDetails(ctx, proj)
}

func (s *service) UpdateProject(ctx context.Context, id uuid.UUID, params UpdateProjectParams) (*entities.ProjectDetails, error) {
	user, err := currentUser(ctx, s.cli)
	if err != nil {
		return nil, err
	}

	proj, err := s.exec.Execute(ctx, id, func(ctx context.Context, cur *entities.Project) ([]events.Event, error) {
		var out []events.Event

		if e, err := events.DecideTitleChanged(
			id, user.ID, cur,
			events.TitleChangedInput{Title: params.Title},
			events.StatusApproved,
		); err != nil {
			return nil, err
		} else {
			emit(cur, &out, e)
		}

		if e, err := events.DecideDescriptionChanged(
			id, user.ID, cur,
			events.DescriptionChangedInput{Description: params.Description},
			events.StatusApproved,
		); err != nil {
			return nil, err
		} else {
			emit(cur, &out, e)
		}

		if e, err := events.DecideStartDateChanged(
			id, user.ID, cur,
			events.StartDateChangedInput{StartDate: params.StartDate},
			events.StatusApproved,
		); err != nil {
			return nil, err
		} else {
			emit(cur, &out, e)
		}

		if e, err := events.DecideEndDateChanged(
			id, user.ID, cur,
			events.EndDateChangedInput{EndDate: params.EndDate},
			events.StatusApproved,
		); err != nil {
			return nil, err
		} else {
			emit(cur, &out, e)
		}

		if e, err := events.DecideOwningOrgNodeChanged(
			id, user.ID, cur,
			events.OwningOrgNodeChangedInput{OwningOrgNodeID: params.OwningOrgNodeID},
			events.StatusApproved,
		); err != nil {
			return nil, err
		} else {
			emit(cur, &out, e)
		}

		return out, nil
	})

	if err != nil {
		return nil, err
	}

	_ = s.cacheProject(ctx, proj)
	return s.buildProjectDetails(ctx, proj)
}

func (s *service) projectRoleID(ctx context.Context, key string) (uuid.UUID, error) {
	var ids []uuid.UUID
	if err := s.cli.ProjectRole.
		Query().
		Where(entprojectrole.KeyEQ(key)).
		Select(entprojectrole.FieldID).
		Scan(ctx, &ids); err != nil {
		return uuid.Nil, err
	}
	if len(ids) == 0 {
		return uuid.Nil, fmt.Errorf("project role not found: %s", key)
	}
	return ids[0], nil
}

// TODO, instead of a helper function there should be a currentUserService
func currentUser(ctx context.Context, cli *ent.Client) (*ent.User, error) {
	authUser, ok := httputil.GetUserFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("no authenticated user in context")
	}

	return cli.User.Get(ctx, authUser.User.ID)
}

func (s *service) GetAllProjects(ctx context.Context) ([]*entities.ProjectDetails, error) {
	user, err := currentUser(ctx, s.cli)
	if err != nil {
		return nil, err
	}

	evts, err := s.cli.Event.
		Query().
		Where(en.TypeEQ(events.ProjectRoleAssignedType)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	uniqueIDs := make(map[uuid.UUID]struct{})
	for _, e := range evts {

		// Safer to marshal/unmarshal to struct
		b, _ := json.Marshal(e.Data)
		var payload *events.ProjectRoleAssigned
		if err := json.Unmarshal(b, &payload); err == nil {
			if payload.PersonID == user.PersonID {
				uniqueIDs[e.ProjectID] = struct{}{}
			}
		}
	}

	projects := make([]*entities.ProjectDetails, 0, len(uniqueIDs))
	for id := range uniqueIDs {
		proj, err := s.fromDb(ctx, id)
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

func (s *service) AddPerson(ctx context.Context, projectID uuid.UUID, personId uuid.UUID) (*entities.ProjectDetails, error) {
	user, err := currentUser(ctx, s.cli)
	if err != nil {
		return nil, err
	}

	contribRoleID, err := s.projectRoleID(ctx, "contributor")
	if err != nil {
		return nil, err
	}

	proj, err := s.exec.Execute(ctx, projectID, func(ctx context.Context, cur *entities.Project) ([]events.Event, error) {
		var out []events.Event
		e, err := events.DecideProjectRoleAssigned(
			projectID, user.ID, cur,
			events.ProjectRoleAssignedInput{PersonID: personId, ProjectRoleID: contribRoleID},
			events.StatusPending,
		)
		if err != nil {
			return nil, err
		}
		emit(cur, &out, e)
		return out, nil
	})
	if err != nil {
		return nil, err
	}

	_ = s.cacheProject(ctx, proj)
	return s.buildProjectDetails(ctx, proj)
}

func (s *service) RemovePerson(ctx context.Context, projectID uuid.UUID, personID uuid.UUID) (*entities.ProjectDetails, error) {
	user, err := currentUser(ctx, s.cli)
	if err != nil {
		return nil, err
	}

	contribRoleID, err := s.projectRoleID(ctx, "contributor")
	if err != nil {
		return nil, err
	}

	proj, err := s.exec.Execute(ctx, projectID, func(ctx context.Context, cur *entities.Project) ([]events.Event, error) {
		var out []events.Event
		e, err := events.DecideProjectRoleUnassigned(
			projectID, user.ID, cur,
			events.ProjectRoleUnassignedInput{PersonID: personID, ProjectRoleID: contribRoleID},
			events.StatusApproved,
		)
		if err != nil {
			return nil, err
		}
		emit(cur, &out, e)
		return out, nil
	})
	if err != nil {
		return nil, err
	}

	_ = s.cacheProject(ctx, proj)
	return s.buildProjectDetails(ctx, proj)
}

func (s *service) UpdateMemberRole(
	ctx context.Context,
	projectID uuid.UUID,
	personID uuid.UUID,
	roleKey string,
) (*entities.ProjectDetails, error) {
	user, err := currentUser(ctx, s.cli)
	if err != nil {
		return nil, err
	}

	newRoleID, err := s.projectRoleID(ctx, roleKey)
	if err != nil {
		return nil, err
	}

	proj, err := s.exec.Execute(ctx, projectID, func(ctx context.Context, cur *entities.Project) ([]events.Event, error) {
		var currentRoleID uuid.UUID
		found := false
		for _, m := range cur.Members {
			if m.PersonID == personID {
				currentRoleID = m.ProjectRoleID
				found = true
				break
			}
		}
		if !found {
			return nil, errors.New("person is not a member of this project")
		}
		if currentRoleID == newRoleID {
			return nil, nil
		}

		var out []events.Event

		unassignEvt, err := events.DecideProjectRoleUnassigned(
			projectID, user.ID, cur,
			events.ProjectRoleUnassignedInput{PersonID: personID, ProjectRoleID: currentRoleID},
			events.StatusApproved,
		)
		if err != nil {
			return nil, err
		}
		emit(cur, &out, unassignEvt)

		assignEvt, err := events.DecideProjectRoleAssigned(
			projectID, user.ID, cur,
			events.ProjectRoleAssignedInput{PersonID: personID, ProjectRoleID: newRoleID},
			events.StatusApproved,
		)
		if err != nil {
			return nil, err
		}
		emit(cur, &out, assignEvt)

		return out, nil
	})
	if err != nil {
		return nil, err
	}

	_ = s.cacheProject(ctx, proj)
	return s.buildProjectDetails(ctx, proj)
}

func (s *service) AddProduct(ctx context.Context, projectID uuid.UUID, productID uuid.UUID) (*entities.ProjectDetails, error) {
	user, err := currentUser(ctx, s.cli)
	if err != nil {
		return nil, err
	}

	proj, err := s.exec.Execute(ctx, projectID, func(ctx context.Context, cur *entities.Project) ([]events.Event, error) {
		var out []events.Event
		e, err := events.DecideProductAdded(
			projectID, user.ID, cur,
			events.ProductAddedInput{ProductID: productID},
			events.StatusApproved,
		)
		if err != nil {
			return nil, err
		}
		emit(cur, &out, e)
		return out, nil
	})
	if err != nil {
		return nil, err
	}

	_ = s.cacheProject(ctx, proj)
	return s.buildProjectDetails(ctx, proj)
}

func (s *service) RemoveProduct(ctx context.Context, projectID uuid.UUID, productID uuid.UUID) (*entities.ProjectDetails, error) {
	user, err := currentUser(ctx, s.cli)
	if err != nil {
		return nil, err
	}

	proj, err := s.exec.Execute(ctx, projectID, func(ctx context.Context, cur *entities.Project) ([]events.Event, error) {
		var out []events.Event
		e, err := events.DecideProductRemoved(
			projectID, user.ID, cur,
			events.ProductRemovedInput{ProductID: productID},
			events.StatusApproved,
		)
		if err != nil {
			return nil, err
		}
		emit(cur, &out, e)
		return out, nil
	})
	if err != nil {
		return nil, err
	}

	_ = s.cacheProject(ctx, proj)
	return s.buildProjectDetails(ctx, proj)
}

func (s *service) buildProjectDetails(ctx context.Context, proj *entities.Project) (*entities.ProjectDetails, error) {
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

	peopleRows, err := s.cli.Person.
		Query().
		Where(personent.IDIn(personIDs...)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	peopleMap := make(map[uuid.UUID]entities.Person)
	for _, p := range peopleRows {
		peopleMap[p.ID] = entities.Person{
			ID:          p.ID,
			UserID:      p.UserID,
			Name:        p.Name,
			GivenName:   p.GivenName,
			FamilyName:  p.FamilyName,
			Email:       p.Email,
			ORCiD:       &p.OrcidID,
			AvatarUrl:   p.AvatarURL,
			Description: p.Description,
		}
	}

	// Fetch Roles
	rolesRows, err := s.cli.ProjectRole.
		Query().
		Where(entprojectrole.IDIn(roleIDs...)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	rolesMap := make(map[uuid.UUID]entities.ProjectRole)
	for _, r := range rolesRows {
		rolesMap[r.ID] = entities.ProjectRole{
			ID:   r.ID,
			Key:  r.Key,
			Name: r.Name,
		}
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

	productRows, err := s.cli.Product.
		Query().
		Where(productent.IDIn(proj.ProductIDs...)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	products := make([]entities.Product, 0, len(productRows))
	for _, p := range productRows {
		products = append(products, entities.Product{
			Id:       p.ID,
			Name:     p.Name,
			Language: *p.Language,
			Type:     entities.ProductType(p.Type),
			DOI:      *p.Doi,
		})
	}

	orgRow, err := s.cli.OrganisationNode.
		Query().
		Where(organisationent.IDEQ(proj.OwningOrgNodeID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	org := entities.OrganisationNode{
		ID:       orgRow.ID,
		ParentID: orgRow.ParentID,
		Name:     orgRow.Name,
	}

	return &entities.ProjectDetails{
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

func (s *service) fromDb(ctx context.Context, projectID uuid.UUID) (*entities.Project, error) {
	// Try cache first
	if proj, err := s.getFromCache(ctx, projectID); err == nil {
		return proj, nil
	}

	evts, version, err := s.es.Load(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if len(evts) == 0 {
		return nil, ErrNotFound
	}

	proj := projection.Reduce(projectID, evts)
	proj.Version = version

	// Update cache
	_ = s.cacheProject(ctx, proj)

	return proj, nil
}

func (s *service) WarmupCache(ctx context.Context) error {
	if s.redis == nil {
		logrus.Warn("Redis not initialized, skipping cache warmup")
	}

	logrus.Info("Starting cache warmup...")

	// Get all project IDs from ProjectStarted events
	var projectIDs []uuid.UUID
	if err := s.cli.Event.Query().
		Where(en.TypeEQ(events.ProjectStartedType)).
		Select(en.FieldProjectID).
		Scan(ctx, &projectIDs); err != nil {
		return err
	}

	count := 0
	for _, id := range projectIDs {

		evts, version, err := s.es.Load(ctx, id)
		if err != nil {
			logrus.Errorf("Failed to load project %s during warmup: %v", id, err)
			continue
		}
		if len(evts) == 0 {
			continue
		}

		proj := projection.Reduce(id, evts)
		proj.Version = version

		if err := s.cacheProject(ctx, proj); err != nil {
			logrus.Errorf("Failed to cache project %s: %v", id, err)
		} else {
			count++
		}
	}

	logrus.Infof("Cache warmup completed. Cached %d projects.", count)
	return nil
}

func (s *service) cacheProject(ctx context.Context, proj *entities.Project) error {
	if s.redis == nil {
		return nil
	}

	bytes, err := json.Marshal(proj)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("project:%s", proj.Id.String())
	return s.redis.Set(ctx, key, bytes, 24*time.Hour).Err()
}

func (s *service) getFromCache(ctx context.Context, projectID uuid.UUID) (*entities.Project, error) {
	if s.redis == nil {
		return nil, errors.New("redis not initialized")
	}

	key := fmt.Sprintf("project:%s", projectID.String())
	val, err := s.redis.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var proj entities.Project
	if err := json.Unmarshal([]byte(val), &proj); err != nil {
		return nil, err
	}

	return &proj, nil
}

func (s *service) GetProjectRoles(ctx context.Context) ([]entities.ProjectRole, error) {
	roles, err := s.cli.ProjectRole.Query().All(ctx)
	if err != nil {
		return nil, err
	}

	var result []entities.ProjectRole
	for _, r := range roles {
		result = append(result, entities.ProjectRole{
			ID:   r.ID,
			Key:  r.Key,
			Name: r.Name,
		})
	}
	return result, nil
}

func (s *service) onStatusChange(ctx context.Context, e events.Event) error {
	projectID := e.AggregateID()

	evts, version, err := s.es.Load(ctx, projectID)
	if err != nil {
		return err
	}
	if len(evts) == 0 {
		return ErrNotFound
	}

	proj := projection.Reduce(projectID, evts)
	proj.Version = version

	// Update cache
	if err := s.cacheProject(ctx, proj); err != nil {
		logrus.Errorf("failed to update project cache on status change: %v", err)
	}

	return nil
}

func emit(cur *entities.Project, out *[]events.Event, e events.Event) {
	if e == nil {
		return
	}
	*out = append(*out, e)
	projection.Apply(cur, e)
}
