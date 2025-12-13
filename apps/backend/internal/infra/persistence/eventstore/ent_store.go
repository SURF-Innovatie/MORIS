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
				SetTitle(v.Title).
				SetDescription(v.Description).
				SetStartDate(v.StartDate).
				SetEndDate(v.EndDate).
				SetOwningOrgNodeID(v.OwningOrgNodeID).
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

		case events.OwningOrgNodeChanged:
			evRow, err := tx.Event.
				Create().
				SetID(id).
				SetProjectID(projectID).
				SetVersion(version).
				SetType(events.OwningOrgNodeChangedType).
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

			if _, err := tx.OwningOrgNodeChangedEvent.
				Create().
				SetEvent(evRow).
				SetOwningOrgNodeID(v.OwningOrgNodeID).
				Save(ctx); err != nil {
				_ = tx.Rollback()
				return err
			}

		case events.ProjectRoleAssigned:
			evRow, err := tx.Event.
				Create().
				SetID(id).
				SetProjectID(projectID).
				SetVersion(version).
				SetType(events.ProjectRoleAssignedType).
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

			if _, err := tx.ProjectRoleAssignedEvent.
				Create().
				SetEvent(evRow).
				SetPersonID(v.PersonID).
				SetProjectRoleID(v.ProjectRoleID).
				Save(ctx); err != nil {
				_ = tx.Rollback()
				return err
			}

		case events.ProjectRoleUnassigned:
			evRow, err := tx.Event.
				Create().
				SetID(id).
				SetProjectID(projectID).
				SetVersion(version).
				SetType(events.ProjectRoleUnassignedType).
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

			if _, err := tx.ProjectRoleUnassignedEvent.
				Create().
				SetEvent(evRow).
				SetPersonID(v.PersonID).
				SetProjectRoleID(v.ProjectRoleID).
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
		WithOwningOrgNodeChanged().
		WithProjectRoleAssigned().
		WithProjectRoleUnassigned().
		WithProductAdded(func(q *ent.ProductAddedEventQuery) {
			q.WithProduct()
		}).
		WithProductRemoved(func(q *ent.ProductRemovedEventQuery) {
			q.WithProduct()
		}).
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
		WithOwningOrgNodeChanged().
		WithProjectRoleAssigned().
		WithProjectRoleUnassigned().
		WithProductAdded(func(q *ent.ProductAddedEventQuery) {
			q.WithProduct()
		}).
		WithProductRemoved(func(q *ent.ProductRemovedEventQuery) {
			q.WithProduct()
		}).
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
			Base:            base,
			Title:           payload.Title,
			Description:     payload.Description,
			StartDate:       payload.StartDate,
			EndDate:         payload.EndDate,
			OwningOrgNodeID: payload.OwningOrgNodeID,
			Members:         nil,
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

	case events.OwningOrgNodeChangedType:
		payload := r.Edges.OwningOrgNodeChanged
		if payload == nil {
			return nil, fmt.Errorf("missing OrganisationChanged edge for event %s", r.ID)
		}
		return events.OwningOrgNodeChanged{
			Base:            base,
			OwningOrgNodeID: payload.OwningOrgNodeID,
		}, nil

	case events.ProjectRoleAssignedType:
		payload := r.Edges.ProjectRoleAssigned
		if payload == nil {
			return nil, fmt.Errorf("missing ProjectRoleAssigned edge for event %s", r.ID)
		}
		return events.ProjectRoleAssigned{
			Base:          base,
			PersonID:      payload.PersonID,
			ProjectRoleID: payload.ProjectRoleID,
		}, nil

	case events.ProjectRoleUnassignedType:
		payload := r.Edges.ProjectRoleUnassigned
		if payload == nil {
			return nil, fmt.Errorf("missing ProjectRoleUnassigned edge for event %s", r.ID)
		}
		return events.ProjectRoleUnassigned{
			Base:          base,
			PersonID:      payload.PersonID,
			ProjectRoleID: payload.ProjectRoleID,
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

func (s *EntStore) LoadUserApprovedEvents(ctx context.Context, userID uuid.UUID) ([]events.Event, error) {
	rows, err := s.cli.Event.
		Query().
		Where(
			en.CreatedByEQ(userID),
			en.StatusEQ(en.StatusApproved),
		).
		Order(ent.Desc(en.FieldOccurredAt)).
		WithProjectStarted().
		WithTitleChanged().
		WithDescriptionChanged().
		WithStartDateChanged().
		WithEndDateChanged().
		WithOwningOrgNodeChanged().
		WithProjectRoleAssigned().
		WithProjectRoleUnassigned().
		WithProductAdded(func(q *ent.ProductAddedEventQuery) {
			q.WithProduct()
		}).
		WithProductRemoved(func(q *ent.ProductRemovedEventQuery) {
			q.WithProduct()
		}).
		All(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]events.Event, 0, len(rows))
	for _, r := range rows {
		evt, err := s.mapEventRow(r)
		if err != nil {
			return nil, err
		}
		if evt != nil {
			out = append(out, evt)
		}
	}

	return out, nil
}
