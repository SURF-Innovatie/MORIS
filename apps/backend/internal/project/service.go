package project

import (
	"context"
	"errors"
	"time"

	personent "github.com/SURF-Innovatie/MORIS/ent/person"
	"github.com/SURF-Innovatie/MORIS/internal/api/organisationdto"
	"github.com/SURF-Innovatie/MORIS/internal/api/persondto"
	"github.com/SURF-Innovatie/MORIS/internal/api/projectdto"
	"github.com/google/uuid"

	"github.com/SURF-Innovatie/MORIS/ent"
	en "github.com/SURF-Innovatie/MORIS/ent/event"
	"github.com/SURF-Innovatie/MORIS/internal/domain/commands"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/SURF-Innovatie/MORIS/internal/domain/projection"
	"github.com/SURF-Innovatie/MORIS/internal/platform/eventstore"
)

// ErrNotFound is returned when a project does not exist (no events).
var ErrNotFound = errors.New("project not found")

type Service interface {
	GetProject(ctx context.Context, id uuid.UUID) (*projectdto.Response, error)
	GetAllProjects(ctx context.Context) ([]*projectdto.Response, error)
	StartProject(ctx context.Context, params StartProjectParams) (*projectdto.Response, error)
	AddPerson(ctx context.Context, projectID uuid.UUID, personID uuid.UUID) (*projectdto.Response, error)
	RemovePerson(ctx context.Context, projectID uuid.UUID, personID uuid.UUID) (*projectdto.Response, error)
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

func (s *service) GetProject(ctx context.Context, id uuid.UUID) (*projectdto.Response, error) {
	evts, version, err := s.es.Load(ctx, id)
	if err != nil {
		return nil, err
	}
	if len(evts) == 0 {
		return nil, ErrNotFound
	}

	proj := projection.Reduce(id, evts)
	proj.Version = version

	resp, err := s.projectToResponse(ctx, proj)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *service) StartProject(ctx context.Context, params StartProjectParams) (*projectdto.Response, error) {
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

	resp, err := s.projectToResponse(ctx, proj)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *service) GetAllProjects(ctx context.Context) ([]*projectdto.Response, error) {
	var ids []uuid.UUID
	if err := s.cli.Event.
		Query().
		Unique(true).
		Select(en.FieldProjectID).
		Scan(ctx, &ids); err != nil {
		return nil, err
	}

	projects := make([]*projectdto.Response, 0, len(ids))
	for _, id := range ids {
		evts, version, err := s.es.Load(ctx, id)
		if err != nil {
			return nil, err
		}
		if len(evts) == 0 {
			continue
		}

		proj := projection.Reduce(id, evts)
		proj.Version = version

		dto, err := s.projectToResponse(ctx, proj)
		if err != nil {
			return nil, err
		}

		projects = append(projects, dto)
	}

	return projects, nil
}

func (s *service) AddPerson(
	ctx context.Context,
	projectID uuid.UUID,
	personId uuid.UUID,
) (*projectdto.Response, error) {
	evts, version, err := s.es.Load(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if len(evts) == 0 {
		return nil, ErrNotFound
	}

	proj := projection.Reduce(projectID, evts)
	proj.Version = version

	evt, err := commands.AddPerson(projectID, proj, personId)
	if err != nil {
		return nil, err
	}
	if evt == nil {
		resp, err := s.projectToResponse(ctx, proj)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}

	if err := s.es.Append(ctx, projectID, version, evt); err != nil {
		return nil, err
	}

	projection.Apply(proj, evt)
	proj.Version = version + 1

	resp, err := s.projectToResponse(ctx, proj)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *service) RemovePerson(
	ctx context.Context,
	projectID uuid.UUID,
	personID uuid.UUID,
) (*projectdto.Response, error) {
	evts, version, err := s.es.Load(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if len(evts) == 0 {
		return nil, ErrNotFound
	}

	proj := projection.Reduce(projectID, evts)
	proj.Version = version

	var personId uuid.UUID
	for _, p := range proj.People {
		if p == personID {
			personId = p
			break
		}
	}

	evt, err := commands.RemovePerson(projectID, proj, personId)
	if err != nil {
		return nil, err
	}
	if evt == nil {
		resp, err := s.projectToResponse(ctx, proj)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}

	if err := s.es.Append(ctx, projectID, version, evt); err != nil {
		return nil, err
	}

	projection.Apply(proj, evt)
	proj.Version = version + 1

	resp, err := s.projectToResponse(ctx, proj)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *service) projectToResponse(ctx context.Context, proj *entities.Project) (*projectdto.Response, error) {
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

	peopleDTOs := make([]persondto.Response, 0, len(peopleRows))
	for _, p := range peopleRows {
		peopleDTOs = append(peopleDTOs, persondto.Response{
			ID:         p.ID,
			Name:       p.Name,
			GivenName:  p.GivenName,
			FamilyName: p.FamilyName,
			Email:      p.Email,
		})
	}

	orgDTO := organisationdto.Response{
		Id:   proj.Organisation.Id,
		Name: proj.Organisation.Name,
	}

	return &projectdto.Response{
		Id:           proj.Id,
		Version:      proj.Version,
		Title:        proj.Title,
		Description:  proj.Description,
		StartDate:    proj.StartDate,
		EndDate:      proj.EndDate,
		Organization: orgDTO,
		People:       peopleDTOs,
	}, nil
}
