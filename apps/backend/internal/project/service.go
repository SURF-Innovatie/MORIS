package project

import (
	"context"
	"errors"

	"github.com/SURF-Innovatie/MORIS/ent"
	en "github.com/SURF-Innovatie/MORIS/ent/event"
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
	AddPerson(ctx context.Context, projectID uuid.UUID, person *entities.Person) (*entities.Project, error)
	RemovePerson(ctx context.Context, projectID uuid.UUID, person *entities.Person) (*entities.Project, error)
}

type service struct {
	cli *ent.Client
	es  eventstore.Store
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

	// Adjust this call if your actual signature differs,
	// but this mirrors AddPerson’s reported signature.
	evt, err := commands.RemovePerson(projectID, proj, *person)
	if err != nil {
		return nil, err
	}
	if evt == nil {
		// No change (e.g. person not present).
		return proj, nil
	}

	if err := s.es.Append(ctx, projectID, version, evt); err != nil {
		return nil, err
	}

	projection.Apply(proj, evt)
	proj.Version = version + 1

	return proj, nil
}
