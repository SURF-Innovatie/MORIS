package project

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"github.com/SURF-Innovatie/MORIS/ent"
	en "github.com/SURF-Innovatie/MORIS/ent/event" //nolint:depguard
	organisationent "github.com/SURF-Innovatie/MORIS/ent/organisationnode"
	personent "github.com/SURF-Innovatie/MORIS/ent/person"
	productent "github.com/SURF-Innovatie/MORIS/ent/product"
	entprojectrole "github.com/SURF-Innovatie/MORIS/ent/projectrole"
	"github.com/SURF-Innovatie/MORIS/internal/app/commandbus"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/SURF-Innovatie/MORIS/internal/domain/projection"
	"github.com/SURF-Innovatie/MORIS/internal/event"
	"github.com/SURF-Innovatie/MORIS/internal/infra/cache"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/eventstore"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// ErrNotFound is returned when a project does not exist (no events).
var ErrNotFound = errors.New("project not found")

type Service interface {
	GetProject(ctx context.Context, id uuid.UUID) (*entities.ProjectDetails, error)
	GetAllProjects(ctx context.Context) ([]*entities.ProjectDetails, error)
	GetChangeLog(ctx context.Context, id uuid.UUID) (*entities.ChangeLog, error)
	GetPendingEvents(ctx context.Context, projectID uuid.UUID) ([]events.DetailedEvent, error)
	GetProjectRoles(ctx context.Context) ([]entities.ProjectRole, error)
	WarmupCache(ctx context.Context) error
}

type service struct {
	cli       *ent.Client
	es        eventstore.Store
	evtSvc    event.Service
	cache     cache.ProjectCache
	refresher cache.ProjectCacheRefresher

	exec *commandbus.Executor[entities.Project]
}

func NewService(es eventstore.Store, cli *ent.Client, evtSvc event.Service, pc cache.ProjectCache, ref cache.ProjectCacheRefresher) Service {
	s := &service{
		es:        es,
		cli:       cli,
		evtSvc:    evtSvc,
		cache:     pc,
		refresher: ref,
	}

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

func (s *service) GetPendingEvents(ctx context.Context, projectID uuid.UUID) ([]events.DetailedEvent, error) {
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

	return s.enrichEvents(ctx, pending)
}

func (s *service) enrichEvents(ctx context.Context, evts []events.Event) ([]events.DetailedEvent, error) {
	if len(evts) == 0 {
		return nil, nil
	}

	// Collect IDs
	personIDs := make(map[uuid.UUID]struct{})
	productIDs := make(map[uuid.UUID]struct{})
	roleIDs := make(map[uuid.UUID]struct{})
	orgIDs := make(map[uuid.UUID]struct{})

	for _, e := range evts {
		if r, ok := e.(events.HasRelatedIDs); ok {
			ids := r.RelatedIDs()
			if ids.PersonID != nil {
				personIDs[*ids.PersonID] = struct{}{}
			}
			if ids.ProductID != nil {
				productIDs[*ids.ProductID] = struct{}{}
			}
			if ids.ProjectRoleID != nil {
				roleIDs[*ids.ProjectRoleID] = struct{}{}
			}
			if ids.OrgNodeID != nil {
				orgIDs[*ids.OrgNodeID] = struct{}{}
			}
		}
	}

	// Fetch Entities
	// Persons
	pIDs := make([]uuid.UUID, 0, len(personIDs))
	for id := range personIDs {
		pIDs = append(pIDs, id)
	}
	personsMap := make(map[uuid.UUID]entities.Person)
	if len(pIDs) > 0 {
		rows, err := s.cli.Person.Query().Where(personent.IDIn(pIDs...)).All(ctx)
		if err != nil {
			return nil, err
		}
		for _, r := range rows {
			personsMap[r.ID] = *(&entities.Person{}).FromEnt(r)
		}
	}

	// Products
	prodIDs := make([]uuid.UUID, 0, len(productIDs))
	for id := range productIDs {
		prodIDs = append(prodIDs, id)
	}
	productsMap := make(map[uuid.UUID]entities.Product)
	if len(prodIDs) > 0 {
		rows, err := s.cli.Product.Query().Where(productent.IDIn(prodIDs...)).All(ctx)
		if err != nil {
			return nil, err
		}
		for _, r := range rows {
			productsMap[r.ID] = entities.Product{
				Id:       r.ID,
				Name:     r.Name,
				Language: *r.Language,
				Type:     entities.ProductType(r.Type),
				DOI:      *r.Doi,
			}
		}
	}

	// Roles
	rIDs := make([]uuid.UUID, 0, len(roleIDs))
	for id := range roleIDs {
		rIDs = append(rIDs, id)
	}
	rolesMap := make(map[uuid.UUID]entities.ProjectRole)
	if len(rIDs) > 0 {
		rows, err := s.cli.ProjectRole.Query().Where(entprojectrole.IDIn(rIDs...)).All(ctx)
		if err != nil {
			return nil, err
		}
		for _, r := range rows {
			rolesMap[r.ID] = entities.ProjectRole{
				ID:   r.ID,
				Key:  r.Key,
				Name: r.Name,
			}
		}
	}

	// Construct DetailedEvents
	result := make([]events.DetailedEvent, len(evts))
	for i, e := range evts {
		de := events.DetailedEvent{Event: e}
		if r, ok := e.(events.HasRelatedIDs); ok {
			ids := r.RelatedIDs()
			if ids.PersonID != nil {
				if p, ok := personsMap[*ids.PersonID]; ok {
					de.Person = &p
				}
			}
			if ids.ProductID != nil {
				if p, ok := productsMap[*ids.ProductID]; ok {
					de.Product = &p
				}
			}
			if ids.ProjectRoleID != nil {
				if r, ok := rolesMap[*ids.ProjectRoleID]; ok {
					de.ProjectRole = &r
				}
			}
		}
		result[i] = de
	}

	return result, nil
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
	if s.cache != nil {
		if proj, err := s.cache.GetProject(ctx, projectID); err == nil {
			return proj, nil
		}
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

	_ = s.cache.SetProject(ctx, proj)
	return proj, nil
}

func (s *service) WarmupCache(ctx context.Context) error {
	if s.cache == nil {
		logrus.Warn("Redis not initialized, skipping cache warmup")
	}

	logrus.Info("Starting cache warmup...")

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

		if err := s.cache.SetProject(ctx, proj); err != nil {
			logrus.Errorf("Failed to cache project %s: %v", id, err)
		} else {
			count++
		}
	}

	logrus.Infof("Cache warmup completed. Cached %d projects.", count)
	return nil
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

	if s.refresher != nil {
		_, err := s.refresher.Refresh(ctx, projectID)
		return err
	}

	evts, version, err := s.es.Load(ctx, projectID)
	if err != nil {
		return err
	}
	if len(evts) == 0 {
		return ErrNotFound
	}

	proj := projection.Reduce(projectID, evts)
	proj.Version = version
	_ = s.cache.SetProject(ctx, proj)
	return nil
}
