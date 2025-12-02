package eventstore

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/SURF-Innovatie/MORIS/ent"
	en "github.com/SURF-Innovatie/MORIS/ent/event"
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

		id := e.GetID()
		if id == uuid.Nil {
			id = uuid.New()
		}

		switch v := e.(type) {
		case events.ProjectStarted:
			// 1) create base Event row
			evRow, err := tx.Event.
				Create().
				SetID(id).
				SetProjectID(projectID).
				SetVersion(version).
				SetType(events.ProjectStartedType).
				SetStatus(en.Status(e.GetStatus())).
				SetOccurredAt(now).
				SetCreatedBy(v.CreatedBy). // Added
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
				SetProjectAdmin(v.ProjectAdmin).
				SetTitle(v.Title).
				SetDescription(v.Description).
				SetStartDate(v.StartDate).
				SetEndDate(v.EndDate).
				SetOrganisationID(v.OrganisationID).
				Save(ctx); err != nil {
				_ = tx.Rollback()
				return err
			}

		case events.TitleChanged:
			evRow, err := tx.Event.
				Create().
				SetID(id).
				SetProjectID(projectID).
				SetVersion(version).
				SetType(events.TitleChangedType).
				SetStatus(en.Status(e.GetStatus())).
				SetOccurredAt(now).
				SetCreatedBy(v.CreatedBy). // Added
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
				SetID(id).
				SetProjectID(projectID).
				SetVersion(version).
				SetType(events.DescriptionChangedType).
				SetStatus(en.Status(e.GetStatus())).
				SetOccurredAt(now).
				SetCreatedBy(v.CreatedBy). // Added
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
				SetID(id).
				SetProjectID(projectID).
				SetVersion(version).
				SetType(events.StartDateChangedType).
				SetOccurredAt(now).
				SetStatus(en.Status(e.GetStatus())).
				SetCreatedBy(v.CreatedBy). // Added
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
				SetID(id).
				SetProjectID(projectID).
				SetVersion(version).
				SetType(events.EndDateChangedType).
				SetOccurredAt(now).
				SetStatus(en.Status(e.GetStatus())).
				SetCreatedBy(v.CreatedBy). // Added
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
				SetID(id).
				SetProjectID(projectID).
				SetVersion(version).
				SetType(events.OrganisationChangedType).
				SetOccurredAt(now).
				SetStatus(en.Status(e.GetStatus())).
				SetCreatedBy(v.CreatedBy). // Added
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
				SetOrganisationID(v.OrganisationID).
				Save(ctx); err != nil {
				_ = tx.Rollback()
				return err
			}

		case events.PersonAdded:
			evRow, err := tx.Event.
				Create().
				SetID(id).
				SetProjectID(projectID).
				SetVersion(version).
				SetType(events.PersonAddedType).
				SetOccurredAt(now).
				SetStatus(en.Status(e.GetStatus())).
				SetCreatedBy(v.CreatedBy). // Added
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
				SetID(id).
				SetProjectID(projectID).
				SetVersion(version).
				SetType(events.PersonRemovedType).
				SetOccurredAt(now).
				SetStatus(en.Status(e.GetStatus())).
				SetCreatedBy(v.CreatedBy). // Added
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
				SetPersonID(v.PersonId).
				Save(ctx); err != nil {
				_ = tx.Rollback()
				return err
			}

		case events.ProductAdded:
			evRow, err := tx.Event.
				Create().
				SetID(id).
				SetProjectID(projectID).
				SetVersion(version).
				SetType(events.ProductAddedType).
				SetOccurredAt(now).
				SetStatus(en.Status(e.GetStatus())).
				SetCreatedBy(v.CreatedBy). // Added
				Save(ctx)
			if err != nil {
				_ = tx.Rollback()
				if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
					return ErrConcurrency
				}
				return err
			}

			if _, err := tx.ProductAddedEvent.
				Create().
				SetEvent(evRow).
				SetProductID(v.ProductID).
				Save(ctx); err != nil {
				_ = tx.Rollback()
				return err
			}

		case events.ProductRemoved:
			evRow, err := tx.Event.
				Create().
				SetID(id).
				SetProjectID(projectID).
				SetVersion(version).
				SetType(events.ProductRemovedType).
				SetOccurredAt(now).
				SetStatus(en.Status(e.GetStatus())).
				SetCreatedBy(v.CreatedBy). // Added
				Save(ctx)
			if err != nil {
				_ = tx.Rollback()
				if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
					return ErrConcurrency
				}
				return err
			}

			if _, err := tx.ProductRemovedEvent.
				Create().
				SetEvent(evRow).
				SetProductID(v.ProductID).
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

func (s *EntStore) UpdateEventStatus(ctx context.Context, eventID uuid.UUID, status string) error {
	// Validate status enum
	switch status {
	case "pending", "approved", "rejected":
	default:
		return fmt.Errorf("invalid status: %s", status)
	}

	return s.cli.Event.
		UpdateOneID(eventID).
		SetStatus(en.Status(status)).
		Exec(ctx)
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
		WithProductAdded().
		WithProductRemoved().
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	out := make([]events.Event, 0, len(rows))

	for _, r := range rows {
		evt, err := s.mapEventRow(r)
		if err != nil {
			return nil, 0, err
		}
		if evt != nil {
			out = append(out, evt)
		}
	}

	curVersion := 0
	if n := len(rows); n > 0 {
		curVersion = rows[n-1].Version
	}
	return out, curVersion, nil
}

func (s *EntStore) LoadEvent(ctx context.Context, eventID uuid.UUID) (events.Event, error) {
	r, err := s.cli.Event.
		Query().
		Where(en.IDEQ(eventID)).
		WithProjectStarted().
		WithTitleChanged().
		WithDescriptionChanged().
		WithStartDateChanged().
		WithEndDateChanged().
		WithOrganisationChanged().
		WithPersonAdded().
		WithPersonRemoved().
		WithProductAdded().
		WithProductRemoved().
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return s.mapEventRow(r)
}

func (s *EntStore) mapEventRow(r *ent.Event) (events.Event, error) {
	base := events.Base{
		ID:        r.ID,
		ProjectID: r.ProjectID,
		At:        r.OccurredAt,
		CreatedBy: r.CreatedBy,
		Status:    string(r.Status),
	}

	switch r.Type {
	case events.ProjectStartedType:
		payload := r.Edges.ProjectStarted
		if payload == nil {
			return nil, fmt.Errorf("missing ProjectStarted edge for event %s", r.ID)
		}
		return events.ProjectStarted{
			Base:           base,
			ProjectAdmin:   payload.ProjectAdmin,
			Title:          payload.Title,
			Description:    payload.Description,
			StartDate:      payload.StartDate,
			EndDate:        payload.EndDate,
			OrganisationID: payload.OrganisationID,
			People:         nil, // driven by PersonAdded events
		}, nil

	case events.TitleChangedType:
		payload := r.Edges.TitleChanged
		if payload == nil {
			return nil, fmt.Errorf("missing TitleChanged edge for event %s", r.ID)
		}
		return events.TitleChanged{
			Base:  base,
			Title: payload.Title,
		}, nil

	case events.DescriptionChangedType:
		payload := r.Edges.DescriptionChanged
		if payload == nil {
			return nil, fmt.Errorf("missing DescriptionChanged edge for event %s", r.ID)
		}
		return events.DescriptionChanged{
			Base:        base,
			Description: payload.Description,
		}, nil

	case events.StartDateChangedType:
		payload := r.Edges.StartDateChanged
		if payload == nil {
			return nil, fmt.Errorf("missing StartDateChanged edge for event %s", r.ID)
		}
		return events.StartDateChanged{
			Base:      base,
			StartDate: payload.StartDate,
		}, nil

	case events.EndDateChangedType:
		payload := r.Edges.EndDateChanged
		if payload == nil {
			return nil, fmt.Errorf("missing EndDateChanged edge for event %s", r.ID)
		}
		return events.EndDateChanged{
			Base:    base,
			EndDate: payload.EndDate,
		}, nil

	case events.OrganisationChangedType:
		payload := r.Edges.OrganisationChanged
		if payload == nil {
			return nil, fmt.Errorf("missing OrganisationChanged edge for event %s", r.ID)
		}
		return events.OrganisationChanged{
			Base:           base,
			OrganisationID: payload.OrganisationID,
		}, nil

	case events.PersonAddedType:
		payload := r.Edges.PersonAdded
		if payload == nil {
			return nil, fmt.Errorf("missing PersonAdded edge for event %s", r.ID)
		}
		return events.PersonAdded{
			Base:     base,
			PersonId: payload.PersonID,
		}, nil

	case events.PersonRemovedType:
		payload := r.Edges.PersonRemoved
		if payload == nil {
			return nil, fmt.Errorf("missing PersonRemoved edge for event %s", r.ID)
		}
		return events.PersonRemoved{
			Base:     base,
			PersonId: payload.PersonID,
		}, nil

	case events.ProductAddedType:
		payload := r.Edges.ProductAdded
		if payload == nil {
			return nil, fmt.Errorf("missing ProductAdded edge for event %s", r.ID)
		}
		return events.ProductAdded{
			Base:      base,
			ProductID: payload.ProductID,
		}, nil

	case events.ProductRemovedType:
		payload := r.Edges.ProductRemoved
		if payload == nil {
			return nil, fmt.Errorf("missing ProductRemoved edge for event %s", r.ID)
		}
		return events.ProductRemoved{
			Base:      base,
			ProductID: payload.ProductID,
		}, nil

	default:
		log.Printf("unknown event type %T", r)
		return nil, nil
	}
}
