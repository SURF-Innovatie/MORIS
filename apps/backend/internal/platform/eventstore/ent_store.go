package eventstore

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/SURF-Innovatie/MORIS/ent"
	en "github.com/SURF-Innovatie/MORIS/ent/event"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
)

type EntStore struct {
	cli *ent.Client
}

func NewEntStore(cli *ent.Client) *EntStore { return &EntStore{cli: cli} }

// Append appends events with optimistic concurrency on (project_id, version).
func (s *EntStore) Append(
	ctx context.Context,
	projectID uuid.UUID,
	expectedVersion int,
	list ...events.Event,
) error {
	if len(list) == 0 {
		return nil
	}

	// Check current version
	last, err := s.cli.Event.
		Query().
		Where(en.ProjectIDEQ(projectID)).
		Order(ent.Desc(en.FieldVersion)).
		First(ctx)
	switch {
	case err == nil:
		if last.Version != expectedVersion {
			return ErrConcurrency
		}
	case ent.IsNotFound(err):
		if expectedVersion != 0 {
			return ErrConcurrency
		}
	default:
		return err
	}

	tx, err := s.cli.Tx(ctx)
	if err != nil {
		return err
	}

	version := expectedVersion
	now := time.Now().UTC()

	for _, e := range list {
		version++

		switch v := e.(type) {

		case events.ProjectStarted:
			// 1) create base Event row
			evRow, err := tx.Event.
				Create().
				SetID(uuid.New()).
				SetProjectID(projectID).
				SetVersion(version).
				SetType(events.ProjectStartedType).
				SetOccurredAt(now).
				Save(ctx)
			if err != nil {
				_ = tx.Rollback()
				if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
					return ErrConcurrency
				}
				return err
			}

			// 2) create concrete payload row, linking back to Event
			if _, err := tx.ProjectStartedEvent.
				Create().
				SetEvent(evRow). // or SetEventID(evRow.ID)
				SetTitle(v.Title).
				SetDescription(v.Description).
				SetStartDate(v.StartDate).
				SetEndDate(v.EndDate).
				SetOrganisationName(v.Organisation.Name).
				Save(ctx); err != nil {
				_ = tx.Rollback()
				return err
			}

		case events.TitleChanged:
			evRow, err := tx.Event.
				Create().
				SetID(uuid.New()).
				SetProjectID(projectID).
				SetVersion(version).
				SetType(events.TitleChangedType).
				SetOccurredAt(now).
				Save(ctx)
			if err != nil {
				_ = tx.Rollback()
				if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
					return ErrConcurrency
				}
				return err
			}

			if _, err := tx.TitleChangedEvent.
				Create().
				SetEvent(evRow).
				SetTitle(v.Title).
				Save(ctx); err != nil {
				_ = tx.Rollback()
				return err
			}

		case events.DescriptionChanged:
			evRow, err := tx.Event.
				Create().
				SetID(uuid.New()).
				SetProjectID(projectID).
				SetVersion(version).
				SetType(events.DescriptionChangedType).
				SetOccurredAt(now).
				Save(ctx)
			if err != nil {
				_ = tx.Rollback()
				if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
					return ErrConcurrency
				}
				return err
			}

			if _, err := tx.DescriptionChangedEvent.
				Create().
				SetEvent(evRow).
				SetDescription(v.Description).
				Save(ctx); err != nil {
				_ = tx.Rollback()
				return err
			}

		case events.StartDateChanged:
			evRow, err := tx.Event.
				Create().
				SetID(uuid.New()).
				SetProjectID(projectID).
				SetVersion(version).
				SetType(events.StartDateChangedType).
				SetOccurredAt(now).
				Save(ctx)
			if err != nil {
				_ = tx.Rollback()
				if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
					return ErrConcurrency
				}
				return err
			}

			if _, err := tx.StartDateChangedEvent.
				Create().
				SetEvent(evRow).
				SetStartDate(v.StartDate).
				Save(ctx); err != nil {
				_ = tx.Rollback()
				return err
			}

		case events.EndDateChanged:
			evRow, err := tx.Event.
				Create().
				SetID(uuid.New()).
				SetProjectID(projectID).
				SetVersion(version).
				SetType(events.EndDateChangedType).
				SetOccurredAt(now).
				Save(ctx)
			if err != nil {
				_ = tx.Rollback()
				if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
					return ErrConcurrency
				}
				return err
			}

			if _, err := tx.EndDateChangedEvent.
				Create().
				SetEvent(evRow).
				SetEndDate(v.EndDate).
				Save(ctx); err != nil {
				_ = tx.Rollback()
				return err
			}

		case events.OrganisationChanged:
			evRow, err := tx.Event.
				Create().
				SetID(uuid.New()).
				SetProjectID(projectID).
				SetVersion(version).
				SetType(events.OrganisationChangedType).
				SetOccurredAt(now).
				Save(ctx)
			if err != nil {
				_ = tx.Rollback()
				if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
					return ErrConcurrency
				}
				return err
			}

			if _, err := tx.OrganisationChangedEvent.
				Create().
				SetEvent(evRow).
				SetOrganisationName(v.Organisation.Name).
				Save(ctx); err != nil {
				_ = tx.Rollback()
				return err
			}

		case events.PersonAdded:
			evRow, err := tx.Event.
				Create().
				SetID(uuid.New()).
				SetProjectID(projectID).
				SetVersion(version).
				SetType(events.PersonAddedType).
				SetOccurredAt(now).
				Save(ctx)
			if err != nil {
				_ = tx.Rollback()
				if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
					return ErrConcurrency
				}
				return err
			}

			if _, err := tx.PersonAddedEvent.
				Create().
				SetEvent(evRow).
				SetPersonID(v.PersonId).
				Save(ctx); err != nil {
				_ = tx.Rollback()
				return err
			}

		case events.PersonRemoved:
			evRow, err := tx.Event.
				Create().
				SetID(uuid.New()).
				SetProjectID(projectID).
				SetVersion(version).
				SetType(events.PersonRemovedType).
				SetOccurredAt(now).
				Save(ctx)
			if err != nil {
				_ = tx.Rollback()
				if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
					return ErrConcurrency
				}
				return err
			}

			if _, err := tx.PersonRemovedEvent.
				Create().
				SetEvent(evRow).
				SetID(v.PersonId).
				Save(ctx); err != nil {
				_ = tx.Rollback()
				return err
			}

		default:
			_ = tx.Rollback()
			return fmt.Errorf("unknown event type %T", e)
		}
	}

	return tx.Commit()
}

func (s *EntStore) Load(
	ctx context.Context,
	projectID uuid.UUID,
) ([]events.Event, int, error) {
	rows, err := s.cli.Event.
		Query().
		Where(en.ProjectIDEQ(projectID)).
		Order(ent.Asc(en.FieldVersion)).
		WithProjectStarted().
		WithTitleChanged().
		WithDescriptionChanged().
		WithStartDateChanged().
		WithEndDateChanged().
		WithOrganisationChanged().
		WithPersonAdded().
		WithPersonRemoved().
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	out := make([]events.Event, 0, len(rows))

	for _, r := range rows {
		base := events.Base{
			ProjectID: projectID,
			At:        r.OccurredAt,
		}

		switch r.Type {
		case events.ProjectStartedType:
			payload := r.Edges.ProjectStarted
			if payload == nil {
				return nil, 0, fmt.Errorf("missing ProjectStarted edge for event %s", r.ID)
			}
			out = append(out, events.ProjectStarted{
				Base:        base,
				Title:       payload.Title,
				Description: payload.Description,
				StartDate:   payload.StartDate,
				EndDate:     payload.EndDate,
				Organisation: entities.Organisation{
					Id:   uuid.Nil, // or real ID if you store it
					Name: payload.OrganisationName,
				},
				People: nil, // driven by PersonAdded events
			})

		case events.TitleChangedType:
			payload := r.Edges.TitleChanged
			if payload == nil {
				return nil, 0, fmt.Errorf("missing TitleChanged edge for event %s", r.ID)
			}
			out = append(out, events.TitleChanged{
				Base:  base,
				Title: payload.Title,
			})

		case events.DescriptionChangedType:
			payload := r.Edges.DescriptionChanged
			if payload == nil {
				return nil, 0, fmt.Errorf("missing DescriptionChanged edge for event %s", r.ID)
			}
			out = append(out, events.DescriptionChanged{
				Base:        base,
				Description: payload.Description,
			})

		case events.StartDateChangedType:
			payload := r.Edges.StartDateChanged
			if payload == nil {
				return nil, 0, fmt.Errorf("missing StartDateChanged edge for event %s", r.ID)
			}
			out = append(out, events.StartDateChanged{
				Base:      base,
				StartDate: payload.StartDate,
			})

		case events.EndDateChangedType:
			payload := r.Edges.EndDateChanged
			if payload == nil {
				return nil, 0, fmt.Errorf("missing EndDateChanged edge for event %s", r.ID)
			}
			out = append(out, events.EndDateChanged{
				Base:    base,
				EndDate: payload.EndDate,
			})

		case events.OrganisationChangedType:
			payload := r.Edges.OrganisationChanged
			if payload == nil {
				return nil, 0, fmt.Errorf("missing OrganisationChanged edge for event %s", r.ID)
			}
			out = append(out, events.OrganisationChanged{
				Base: base,
				Organisation: entities.Organisation{
					Id:   uuid.Nil,
					Name: payload.OrganisationName,
				},
			})

		case events.PersonAddedType:
			payload := r.Edges.PersonAdded
			if payload == nil {
				return nil, 0, fmt.Errorf("missing PersonAdded edge for event %s", r.ID)
			}
			out = append(out, events.PersonAdded{
				Base:     base,
				PersonId: payload.PersonID,
			})

		case events.PersonRemovedType:
			payload := r.Edges.PersonRemoved
			if payload == nil {
				return nil, 0, fmt.Errorf("missing PersonRemoved edge for event %s", r.ID)
			}
			out = append(out, events.PersonRemoved{
				Base:     base,
				PersonId: payload.ID,
			})

		default:
			// unknown type: ignore or log, depending on your policy
		}
	}

	curVersion := 0
	if n := len(rows); n > 0 {
		curVersion = rows[n-1].Version
	}
	return out, curVersion, nil
}
