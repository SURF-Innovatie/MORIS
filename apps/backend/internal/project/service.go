package project

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/SURF-Innovatie/MORIS/ent"
	en "github.com/SURF-Innovatie/MORIS/ent/event"
	organisationent "github.com/SURF-Innovatie/MORIS/ent/organisation"
	personent "github.com/SURF-Innovatie/MORIS/ent/person"
	"github.com/SURF-Innovatie/MORIS/ent/personaddedevent"
	productent "github.com/SURF-Innovatie/MORIS/ent/product"
	"github.com/SURF-Innovatie/MORIS/ent/projectstartedevent"
	"github.com/SURF-Innovatie/MORIS/internal/domain/commands"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/SURF-Innovatie/MORIS/internal/domain/projection"
	"github.com/SURF-Innovatie/MORIS/internal/event"
	"github.com/SURF-Innovatie/MORIS/internal/handler/middleware"
	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/eventstore"
	"github.com/google/uuid"
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
	GetChangeLog(ctx context.Context, id uuid.UUID) (*entities.ChangeLog, error)
	GetPendingEvents(ctx context.Context, projectID uuid.UUID) ([]events.Event, error)
}

type service struct {
	cli    *ent.Client
	es     eventstore.Store
	evtSvc event.Service
}

type StartProjectParams struct {
	ProjectAdmin   uuid.UUID
	Title          string
	Description    string
	OrganisationID uuid.UUID
	StartDate      time.Time
	EndDate        time.Time
}

type UpdateProjectParams struct {
	Title          string
	Description    string
	OrganisationID uuid.UUID
	StartDate      time.Time
	EndDate        time.Time
}

func NewService(es eventstore.Store, cli *ent.Client, evtSvc event.Service) Service {
	return &service{es: es, cli: cli, evtSvc: evtSvc}
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

	startEvent := events.ProjectStarted{
		Base: events.Base{
			ID:        uuid.New(),
			ProjectID: projectID,
			At:        now,
			Status:    string(en.StatusApproved),
			CreatedBy: user.ID,
		},
		ProjectAdmin:   user.PersonID,
		Title:          params.Title,
		Description:    params.Description,
		StartDate:      params.StartDate,
		EndDate:        params.EndDate,
		OrganisationID: params.OrganisationID,
	}

	if err := s.es.Append(ctx, projectID, 0, startEvent); err != nil {
		return nil, err
	}

	proj := projection.Reduce(projectID, []events.Event{startEvent})
	proj.Version = 1

	ev, err := commands.AddPerson(projectID, user.ID, proj, user.PersonID, en.StatusApproved)
	if err != nil {
		return nil, err
	}

	if err := s.es.Append(ctx, projectID, proj.Version, ev); err != nil {
		return nil, err
	}

	projection.Apply(proj, ev)

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

	if evt, err := commands.SetOrganisation(id, user.ID, proj, params.OrganisationID, en.StatusApproved); err != nil {
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

	_ = s.evtSvc.HandleEvents(ctx, newEvents...)

	return s.buildProjectDetails(ctx, proj)
}

// TODO, instead of a helper function there should be a currentUserService
func currentUser(ctx context.Context, cli *ent.Client) (*ent.User, error) {
	authUser, ok := middleware.GetUserFromContext(ctx)
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

	// 1. Projects where user is admin (ProjectStartedEvent)
	var adminProjectIDs []uuid.UUID
	if err := s.cli.ProjectStartedEvent.
		Query().
		Where(projectstartedevent.ProjectAdminEQ(user.PersonID)).
		QueryEvent().
		Select(en.FieldProjectID).
		Scan(ctx, &adminProjectIDs); err != nil {
		return nil, err
	}

	// 2. Projects where user is added (PersonAddedEvent)
	var memberProjectIDs []uuid.UUID
	if err := s.cli.PersonAddedEvent.
		Query().
		Where(personaddedevent.PersonIDEQ(user.PersonID)).
		QueryEvent().
		Select(en.FieldProjectID).
		Scan(ctx, &memberProjectIDs); err != nil {
		return nil, err
	}

	// Combine IDs
	uniqueIDs := make(map[uuid.UUID]struct{})
	for _, id := range adminProjectIDs {
		uniqueIDs[id] = struct{}{}
	}
	for _, id := range memberProjectIDs {
		uniqueIDs[id] = struct{}{}
	}

	projects := make([]*entities.ProjectDetails, 0, len(uniqueIDs))
	for id := range uniqueIDs {
		proj, err := s.fromDb(ctx, id)
		if err != nil {
			// Skip if not found (shouldn't happen if consistent)
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

func (s *service) AddPerson(
	ctx context.Context,
	projectID uuid.UUID,
	personId uuid.UUID,
) (*entities.ProjectDetails, error) {
	proj, err := s.fromDb(ctx, projectID)
	if err != nil {
		return nil, err
	}

	user, err := currentUser(ctx, s.cli)
	if err != nil {
		return nil, err
	}

	evt, err := commands.AddPerson(projectID, user.ID, proj, personId, en.StatusPending)
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
	proj.Version++

	resp, err := s.buildProjectDetails(ctx, proj)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *service) RemovePerson(
	ctx context.Context,
	projectID uuid.UUID,
	personID uuid.UUID,
) (*entities.ProjectDetails, error) {
	proj, err := s.fromDb(ctx, projectID)
	if err != nil {
		return nil, err
	}

	user, err := currentUser(ctx, s.cli)
	if err != nil {
		return nil, err
	}

	evt, err := commands.RemovePerson(projectID, user.ID, proj, personID, en.StatusApproved)
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
	proj.Version++

	resp, err := s.buildProjectDetails(ctx, proj)
	if err != nil {
		return nil, err
	}

	return resp, nil
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

	peopleRows, err := s.cli.Person.
		Query().
		Where(personent.IDIn(proj.People...)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	people := make([]entities.Person, 0, len(peopleRows))
	for _, p := range peopleRows {
		people = append(people, entities.Person{
			Id:          p.ID,
			UserID:      p.UserID,
			Name:        p.Name,
			GivenName:   p.GivenName,
			FamilyName:  p.FamilyName,
			Email:       p.Email,
			ORCiD:       &p.OrcidID,
			AvatarUrl:   p.AvatarURL,
			Description: p.Description,
		})
	}

	productRows, err := s.cli.Product.
		Query().
		Where(productent.IDIn(proj.Products...)).
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

	orgRow, err := s.cli.Organisation.
		Query().
		Where(organisationent.ID(proj.Organisation)).
		First(ctx)
	if err != nil {
		return nil, err
	}

	org := entities.Organisation{
		Id:   orgRow.ID,
		Name: orgRow.Name,
	}

	return &entities.ProjectDetails{
		Project:      *proj,
		Organisation: org,
		People:       people,
		Products:     products,
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
	evts, version, err := s.es.Load(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if len(evts) == 0 {
		return nil, ErrNotFound
	}

	proj := projection.Reduce(projectID, evts)
	proj.Version = version

	return proj, nil
}
