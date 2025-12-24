package eventstore

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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

	builders := make([]*ent.EventCreate, len(list))
	for i, e := range list {
		version++

		id := e.GetID()
		if id == uuid.Nil {
			id = uuid.New()
		}

		var createdBy uuid.UUID
		if cb, ok := any(e).(interface{ CreatedByID() uuid.UUID }); ok {
			createdBy = cb.CreatedByID()
		}

		dataMap, err := eventToMap(e)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to marshal event %T: %w", e, err)
		}

		builders[i] = tx.Event.
			Create().
			SetID(id).
			SetProjectID(projectID).
			SetVersion(version).
			SetType(e.Type()).
			SetStatus(en.Status(e.GetStatus())).
			SetOccurredAt(now).
			SetCreatedBy(createdBy).
			SetData(dataMap)
	}

	if err := tx.Event.CreateBulk(builders...).Exec(ctx); err != nil {
		_ = tx.Rollback()
		if ent.IsConstraintError(err) {
			return ErrConcurrency
		}
		return err
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
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	out := make([]events.Event, 0, len(rows))
	for _, r := range rows {
		evt, err := s.mapEventRow(r)
		if err != nil {
			// Log error but maybe continue? Critical data corruption if event is unreadable.
			// Return error for now.
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
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return s.mapEventRow(r)
}

func (s *EntStore) LoadUserApprovedEvents(ctx context.Context, userID uuid.UUID) ([]events.Event, error) {
	rows, err := s.cli.Event.
		Query().
		Where(
			en.CreatedByEQ(userID),
			en.StatusEQ(en.StatusApproved),
		).
		Order(ent.Desc(en.FieldOccurredAt)).
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

func (s *EntStore) mapEventRow(r *ent.Event) (events.Event, error) {
	base := events.Base{
		ID:        r.ID,
		ProjectID: r.ProjectID,
		At:        r.OccurredAt,
		CreatedBy: r.CreatedBy,
		Status:    string(r.Status),
	}

	evt, err := events.Create(r.Type)
	if err != nil {
		log.Printf("unknown event type %s", r.Type)
		return nil, nil
	}

	b, err := json.Marshal(r.Data)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, evt); err != nil {
		return nil, err
	}

	evt.SetBase(base)

	return evt, nil
}

func eventToMap(e events.Event) (map[string]any, error) {
	b, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	err = json.Unmarshal(b, &m)
	return m, err
}
