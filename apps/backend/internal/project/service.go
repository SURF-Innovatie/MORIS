package project

import (
	"context"
	"errors"
	"time"

	"github.com/SURF-Innovatie/MORIS/ent"
	en "github.com/SURF-Innovatie/MORIS/ent/event"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/google/uuid"

	"github.com/SURF-Innovatie/MORIS/internal/domain/commands"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/projection"
	"github.com/SURF-Innovatie/MORIS/internal/platform/eventstore"
)

// ErrNotFound is returned when a project does not exist (no events).
var ErrNotFound = errors.New("project not found")

type Service interface {
	GetProject(ctx context.Context, id uuid.UUID) (*entities.Project, error)
	GetAllProjects(ctx context.Context) ([]entities.Project, error)
	StartProject(ctx context.Context, params StartProjectParams) (*entities.Project, error)
	AddPerson(ctx context.Context, projectID uuid.UUID, person *entities.Person) (*entities.Project, error)
	RemovePerson(ctx context.Context, projectID uuid.UUID, person uuid.UUID) (*entities.Project, error)
}

type service struct {
	cli *ent.Client
	es  eventstore.Store
}

type StartProjectParams struct {
	Title        string
	Description  string
	Organisation string
	StartDate    time.Time
	EndDate      time.Time
}

func NewService(es eventstore.Store, cli *ent.Client) Service {
	return &service{es: es, cli: cli}
}

func (s *service) GetProject(ctx context.Context, id uuid.UUID) (*entities.Project, error) {
	evts, version, err := s.es.Load(ctx, id)
	if err != nil {
		return nil, err
	}
	if len(evts) == 0 {
		return nil, ErrNotFound
	}

	proj := projection.Reduce(id, evts)
	proj.Version = version
	return proj, nil
}

func (s *service) StartProject(ctx context.Context, params StartProjectParams) (*entities.Project, error) {
	if params.Title == "" {
		return nil, errors.New("title is required")
	}
	projectID := uuid.New()
	now := time.Now().UTC()

	org := entities.Organisation{
		Id:   uuid.New(),
		Name: params.Organisation,
	}

	startEvent := events.ProjectStarted{
		Base: events.Base{
			ProjectID: projectID,
			At:        now,
		},
		Title:        params.Title,
		Description:  params.Description,
		StartDate:    params.StartDate,
		EndDate:      params.EndDate,
		Organisation: org,
	}

	if err := s.es.Append(ctx, projectID, 0, startEvent); err != nil {
		return nil, err
	}

	proj := projection.Reduce(projectID, []events.Event{startEvent})
	proj.Version = 1

	return proj, nil
}

func (s *service) GetAllProjects(ctx context.Context) ([]entities.Project, error) {
	var ids []uuid.UUID
	if err := s.cli.Event.
		Query().
		Unique(true).
		Select(en.FieldProjectID).
		Scan(ctx, &ids); err != nil {
		return nil, err
	}

	projects := make([]entities.Project, 0, len(ids))
	for _, id := range ids {
		evts, version, err := s.es.Load(ctx, id)
		if err != nil {
			return nil, err
		}
		if len(evts) == 0 {
			return nil, ErrNotFound
		}
		proj := projection.Reduce(id, evts)
		proj.Version = version
		projects = append(projects, *proj)
	}

	return projects, nil
}

func (s *service) AddPerson(
	ctx context.Context,
	projectID uuid.UUID,
	person *entities.Person,
) (*entities.Project, error) {
	evts, version, err := s.es.Load(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if len(evts) == 0 {
		return nil, ErrNotFound
	}

	proj := projection.Reduce(projectID, evts)
	proj.Version = version

	if person == nil {
		return nil, errors.New("person is nil")
	}

	// Domain command decides if an event is needed and enforces rules.
	evt, err := commands.AddPerson(projectID, proj, *person)
	if err != nil {
		return nil, err
	}
	if evt == nil {
		// No change (e.g. person already present).
		return proj, nil
	}

	if err := s.es.Append(ctx, projectID, version, evt); err != nil {
		// Here you could special-case eventstore.ErrConcurrency if you want.
		return nil, err
	}

	// Update in-memory projection with the new event
	projection.Apply(proj, evt)
	proj.Version = version + 1

	return proj, nil
}

func (s *service) RemovePerson(
	ctx context.Context,
	projectID uuid.UUID,
	personId uuid.UUID,
) (*entities.Project, error) {
	evts, version, err := s.es.Load(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if len(evts) == 0 {
		return nil, ErrNotFound
	}

	proj := projection.Reduce(projectID, evts)
	proj.Version = version

	var person *entities.Person

	for _, p := range proj.People {
		if p.Id == personId {
			person = p
			break
		}
	}

	if person == nil {
		return nil, ErrNotFound
	}

	evt, err := commands.RemovePerson(projectID, proj, *person)
	if err != nil {
		return nil, err
	}
	if evt == nil {
		return proj, nil
	}

	if err := s.es.Append(ctx, projectID, version, evt); err != nil {
		return nil, err
	}

	projection.Apply(proj, evt)
	proj.Version = version + 1

	return proj, nil
}
