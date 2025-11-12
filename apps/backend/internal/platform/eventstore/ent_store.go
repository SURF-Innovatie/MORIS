package eventstore

import (
	"context"
	"encoding/json"
	"errors"
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
// expectedVersion is the current stream version before appending.
// Pass 0 for a new stream (first appended event will be version 1).
func (s *EntStore) Append(ctx context.Context, projectID uuid.UUID, expectedVersion int, list []events.Event) error {
	if len(list) == 0 {
		return nil
	}
	// Check current version.
	last, err := s.cli.Event.
		Query().
		Where(en.ProjectIDEQ(projectID)).
		Order(ent.Desc(en.FieldVersion)).
		First(ctx)
	switch {
	case err == nil:
		if last.Version != expectedVersion {
			return errors.New("concurrency conflict")
		}
	case ent.IsNotFound(err):
		if expectedVersion != 0 {
			return errors.New("concurrency conflict")
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
		typ, payload, err := marshal(e)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		if _, err := tx.Event.Create().
			SetID(uuid.New()).
			SetProjectID(projectID).
			SetVersion(version).
			SetType(typ).
			SetData(payload).
			SetOccurredAt(now).
			Save(ctx); err != nil {
			_ = tx.Rollback()
			// translate unique violation to concurrency error
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				return errors.New("concurrency conflict")
			}
			return err
		}
	}
	return tx.Commit()
}

func (s *EntStore) Load(ctx context.Context, projectID uuid.UUID) ([]events.Event, int, error) {
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
		ev, err := unmarshal(r.Type, r.Data, projectID, r.OccurredAt)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, ev)
	}
	curVersion := 0
	if n := len(rows); n > 0 {
		curVersion = rows[n-1].Version
	}
	return out, curVersion, nil
}

// ── JSON codec bound to your events package ───────────────────────────────────

func marshal(e events.Event) (string, []byte, error) {
	switch v := e.(type) {
	case events.ProjectStarted:
		b, err := json.Marshal(v)
		return events.ProjectStartedType, b, err
	case events.TitleChanged:
		b, err := json.Marshal(v)
		return events.TitleChangedType, b, err
	case events.DescriptionChanged:
		b, err := json.Marshal(v)
		return events.DescriptionChangedType, b, err
	case events.StartDateChanged:
		b, err := json.Marshal(v)
		return events.StartDateChangedType, b, err
	case events.EndDateChanged:
		b, err := json.Marshal(v)
		return events.EndDateChangedType, b, err
	case events.OrganisationChanged:
		b, err := json.Marshal(v)
		return events.OrganisationChangedType, b, err
	case events.PersonAdded:
		b, err := json.Marshal(v)
		return events.PersonAddedType, b, err
	case events.PersonRemoved:
		b, err := json.Marshal(v)
		return events.PersonRemovedType, b, err
	default:
		return "", nil, errors.New("unknown event type")
	}
}

func unmarshal(typ string, data []byte, projectID uuid.UUID, at time.Time) (events.Event, error) {
	base := events.Base{ProjectID: projectID, At: at}
	switch typ {
	case events.ProjectStartedType:
		var v events.ProjectStarted
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}
		if v.ProjectID == uuid.Nil {
			v.Base = base
		}
		return v, nil
	case events.TitleChangedType:
		var v events.TitleChanged
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}
		if v.ProjectID == uuid.Nil {
			v.Base = base
		}
		return v, nil
	case events.DescriptionChangedType:
		var v events.DescriptionChanged
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}
		if v.ProjectID == uuid.Nil {
			v.Base = base
		}
		return v, nil
	case events.StartDateChangedType:
		var v events.StartDateChanged
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}
		if v.ProjectID == uuid.Nil {
			v.Base = base
		}
		return v, nil
	case events.EndDateChangedType:
		var v events.EndDateChanged
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}
		if v.ProjectID == uuid.Nil {
			v.Base = base
		}
	case events.OrganisationChangedType:
		var v events.OrganisationChanged
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}
		if v.ProjectID == uuid.Nil {
			v.Base = base
		}
		return v, nil
	case events.PersonAddedType:
		var v events.PersonAdded
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}
		if v.ProjectID == uuid.Nil {
			v.Base = base
		}
		return v, nil
	case events.PersonRemovedType:
		var v events.PersonRemoved
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}
		if v.ProjectID == uuid.Nil {
			v.Base = base
		}
		return v, nil
	default:
		return nil, errors.New("unknown event type: " + typ)
	}
}
