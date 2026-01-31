package event_test

import (
	"context"
	"testing"
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/event"
	"github.com/google/uuid"

	"github.com/SURF-Innovatie/MORIS/ent/enttest"
	_ "github.com/mattn/go-sqlite3"
)

func TestEntStore_AppendAndLoad(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	store := event.NewEntRepo(client)
	ctx := context.Background()

	projectID := uuid.New()
	personID := uuid.New()
	now := time.Now().UTC().Truncate(time.Millisecond) // Truncate for DB precision

	// Create a sample event
	evt := &events.ProjectStarted{
		Base: events.Base{
			ID:        uuid.New(),
			ProjectID: projectID,
			At:        now,
			Status:    "approved",
			CreatedBy: personID,
		},
		Title:           "Test Project",
		Description:     "A test description",
		StartDate:       now,
		EndDate:         now.Add(24 * time.Hour),
		OwningOrgNodeID: uuid.New(),
	}

	// Append
	if err := store.Append(ctx, projectID, 0, evt); err != nil {
		t.Fatalf("failed to append event: %v", err)
	}

	// Load
	loadedEvents, version, err := store.Load(ctx, projectID)
	if err != nil {
		t.Fatalf("failed to load events: %v", err)
	}

	if len(loadedEvents) != 1 {
		t.Fatalf("expected 1 event, got %d", len(loadedEvents))
	}

	if version != 1 {
		t.Fatalf("expected version 1, got %d", version)
	}

	loadedEvt, ok := loadedEvents[0].(*events.ProjectStarted)
	if !ok {
		t.Fatalf("expected *ProjectStarted event, got %T", loadedEvents[0])
	}

	if loadedEvt.Title != evt.Title {
		t.Errorf("expected title %q, got %q", evt.Title, loadedEvt.Title)
	}
	if loadedEvt.Description != evt.Description {
		t.Errorf("expected description %q, got %q", evt.Description, loadedEvt.Description)
	}
}
