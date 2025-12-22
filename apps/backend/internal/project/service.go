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
	"github.com/SURF-Innovatie/MORIS/internal/domain/commands"
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
	now := time.Now().UTC()

	user, err := currentUser(ctx, s.cli)
	if err != nil {
		return nil, err
	}

	startEvent := &events.ProjectStarted{
		Base: events.Base{
			ID:        uuid.New(),
			ProjectID: projectID,
			At:        now,
			Status:    string(en.StatusApproved),
			CreatedBy: user.ID,
		},
		Title:           params.Title,
		Description:     params.Description,
		StartDate:       params.StartDate,
		EndDate:         params.EndDate,
		OwningOrgNodeID: params.OwningOrgNodeID,
	}

	if err := s.es.Append(ctx, projectID, 0, startEvent); err != nil {
		return nil, err
	}

	proj := projection.Reduce(projectID, []events.Event{startEvent})
	proj.Version = 1

	adminRoleID, err := s.projectRoleID(ctx, "admin") // or "lead" depending on your key
	if err != nil {
		return nil, err
	}

	ev, err := commands.AssignProjectRole(
		projectID,
		user.ID,
		proj,
		user.PersonID,
		adminRoleID,
		en.StatusApproved,
	)
	if err != nil {
		return nil, err
	}

	if err := s.es.Append(ctx, projectID, proj.Version, ev); err != nil {
		return nil, err
	}

	projection.Apply(proj, ev)

	// Cache the project
	if err := s.cacheProject(ctx, proj); err != nil {
		// Log error but continue, caching failure shouldn't stop flow
		logrus.Errorf("failed to cache project: %v", err)
	}

	// user is already fetched at the beginning of the function
	_ = s.evtSvc.HandleEvents(ctx, startEvent)

	resp, err := s.buildProjectDetails(ctx, proj)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *service) UpdateProject(ctx context.Context, id uuid.UUID, params UpdateProjectParams) (*entities.ProjectDetails, error) {
	proj, err := s.fromDb(ctx, id)
	if err != nil {
		return nil, err
	}

	user, err := currentUser(ctx, s.cli)
	if err != nil {
		return nil, err
	}

	var newEvents []events.Event

	if evt, err := commands.ChangeTitle(id, user.ID, proj, params.Title, en.StatusApproved); err != nil {
		return nil, err
	} else if evt != nil {
		newEvents = append(newEvents, evt)
		projection.Apply(proj, evt)
	}

	if evt, err := commands.ChangeDescription(id, user.ID, proj, params.Description, en.StatusApproved); err != nil {
		return nil, err
	} else if evt != nil {
		newEvents = append(newEvents, evt)
		projection.Apply(proj, evt)
	}

	if evt, err := commands.ChangeStartDate(id, user.ID, proj, params.StartDate, en.StatusApproved); err != nil {
		return nil, err
	} else if evt != nil {
		newEvents = append(newEvents, evt)
		projection.Apply(proj, evt)
	}

	if evt, err := commands.ChangeEndDate(id, user.ID, proj, params.EndDate, en.StatusApproved); err != nil {
		return nil, err
	} else if evt != nil {
		newEvents = append(newEvents, evt)
		projection.Apply(proj, evt)
	}

	if evt, err := commands.ChangeOwningOrgNode(id, user.ID, proj, params.OwningOrgNodeID, en.StatusApproved); err != nil {
		return nil, err
	} else if evt != nil {
		newEvents = append(newEvents, evt)
		projection.Apply(proj, evt)
	}

	if len(newEvents) == 0 {
		return s.buildProjectDetails(ctx, proj)
	}

	for _, evt := range newEvents {
		if err := s.es.Append(ctx, id, proj.Version, evt); err != nil {
			return nil, err
		}
		proj.Version++
	}

	// Update cache
	if err := s.cacheProject(ctx, proj); err != nil {
		logrus.Errorf("failed to update project cache: %v\n", err)
	}

	_ = s.evtSvc.HandleEvents(ctx, newEvents...)

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
	proj, err := s.fromDb(ctx, projectID)
	if err != nil {
		return nil, err
	}

	user, err := currentUser(ctx, s.cli)
	if err != nil {
		return nil, err
	}

	contribRoleID, err := s.projectRoleID(ctx, "contributor")
	if err != nil {
		return nil, err
	}

	evt, err := commands.AssignProjectRole(projectID, user.ID, proj, personId, contribRoleID, en.StatusPending)
	if err != nil {
		return nil, err
	}
	if evt == nil {
		return s.buildProjectDetails(ctx, proj)
	}

	if err := s.es.Append(ctx, projectID, proj.Version, evt); err != nil {
		return nil, err
	}

	_ = s.evtSvc.HandleEvents(ctx, evt)

	projection.Apply(proj, evt)
	proj.Version++

	_ = s.cacheProject(ctx, proj)
	return s.buildProjectDetails(ctx, proj)
}

func (s *service) RemovePerson(ctx context.Context, projectID uuid.UUID, personID uuid.UUID) (*entities.ProjectDetails, error) {
	proj, err := s.fromDb(ctx, projectID)
	if err != nil {
		return nil, err
	}

	user, err := currentUser(ctx, s.cli)
	if err != nil {
		return nil, err
	}

	contribRoleID, err := s.projectRoleID(ctx, "contributor")
	if err != nil {
		return nil, err
	}

	evt, err := commands.UnassignProjectRole(projectID, user.ID, proj, personID, contribRoleID, en.StatusApproved)
	if err != nil {
		return nil, err
	}
	if evt == nil {
		return s.buildProjectDetails(ctx, proj)
	}

	if err := s.es.Append(ctx, projectID, proj.Version, evt); err != nil {
		return nil, err
	}

	_ = s.evtSvc.HandleEvents(ctx, evt)

	projection.Apply(proj, evt)
	proj.Version++

	_ = s.cacheProject(ctx, proj)
	return s.buildProjectDetails(ctx, proj)
}

func (s *service) UpdateMemberRole(
	ctx context.Context,
	projectID uuid.UUID,
	personID uuid.UUID,
	roleKey string,
) (*entities.ProjectDetails, error) {
	proj, err := s.fromDb(ctx, projectID)
	if err != nil {
		return nil, err
	}

	user, err := currentUser(ctx, s.cli)
	if err != nil {
		return nil, err
	}

	// 1. Find current role
	var currentRoleID uuid.UUID
	found := false
	for _, m := range proj.Members {
		if m.PersonID == personID {
			currentRoleID = m.ProjectRoleID
			found = true
			break
		}
	}
	if !found {
		return nil, errors.New("person is not a member of this project")
	}

	// 2. Resolve new role ID
	newRoleID, err := s.projectRoleID(ctx, roleKey)
	if err != nil {
		return nil, err
	}

	if currentRoleID == newRoleID {
		return s.buildProjectDetails(ctx, proj)
	}

	// 3. Create Events
	// Unassign old role
	unassignEvt, err := commands.UnassignProjectRole(projectID, user.ID, proj, personID, currentRoleID, en.StatusApproved)
	if err != nil {
		return nil, err
	}

	// Assign new role
	assignEvt, err := commands.AssignProjectRole(projectID, user.ID, proj, personID, newRoleID, en.StatusApproved)
	if err != nil {
		return nil, err
	}

	newEvents := []events.Event{}
	if unassignEvt != nil {
		newEvents = append(newEvents, unassignEvt)
	}
	if assignEvt != nil {
		newEvents = append(newEvents, assignEvt)
	}

	if len(newEvents) == 0 {
		return s.buildProjectDetails(ctx, proj)
	}

	// 4. Append and Apply
	for _, evt := range newEvents {
		if err := s.es.Append(ctx, projectID, proj.Version, evt); err != nil {
			return nil, err
		}
		projection.Apply(proj, evt)
		proj.Version++
	}

	_ = s.evtSvc.HandleEvents(ctx, newEvents...)
	_ = s.cacheProject(ctx, proj)

	return s.buildProjectDetails(ctx, proj)
}

func (s *service) AddProduct(
	ctx context.Context,
	projectID uuid.UUID,
	productID uuid.UUID,
) (*entities.ProjectDetails, error) {
	proj, err := s.fromDb(ctx, projectID)
	if err != nil {
		return nil, err
	}

	user, err := currentUser(ctx, s.cli)
	if err != nil {
		return nil, err
	}

	evt, err := commands.AddProduct(projectID, user.ID, proj, productID, en.StatusApproved)
	if err != nil {
		return nil, err
	}
	if evt == nil {
		resp, err := s.buildProjectDetails(ctx, proj)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}

	if err := s.es.Append(ctx, projectID, proj.Version, evt); err != nil {
		return nil, err
	}

	projection.Apply(proj, evt)
	proj.Version += 1

	// Update cache
	if err := s.cacheProject(ctx, proj); err != nil {
		logrus.Errorf("failed to update project cache: %v", err)
	}

	// notify all users about the new product (temporary for demo)
	// Handled by EventService now
	_ = s.evtSvc.HandleEvents(ctx, evt)

	resp, err := s.buildProjectDetails(ctx, proj)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *service) RemoveProduct(
	ctx context.Context,
	projectID uuid.UUID,
	productID uuid.UUID,
) (*entities.ProjectDetails, error) {
	proj, err := s.fromDb(ctx, projectID)
	if err != nil {
		return nil, err
	}

	user, err := currentUser(ctx, s.cli)
	if err != nil {
		return nil, err
	}

	evt, err := commands.RemoveProduct(projectID, user.ID, proj, productID, en.StatusApproved)
	if err != nil {
		return nil, err
	}
	if evt == nil {
		resp, err := s.buildProjectDetails(ctx, proj)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}

	if err := s.es.Append(ctx, projectID, proj.Version, evt); err != nil {
		return nil, err
	}

	_ = s.evtSvc.HandleEvents(ctx, evt)

	projection.Apply(proj, evt)
	proj.Version += 1

	// Update cache
	if err := s.cacheProject(ctx, proj); err != nil {
		logrus.Errorf("failed to update project cache: %v", err)
	}

	resp, err := s.buildProjectDetails(ctx, proj)
	if err != nil {
		return nil, err
	}

	return resp, nil
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
