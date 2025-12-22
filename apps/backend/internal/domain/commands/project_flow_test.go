package commands_test

import (
	"testing"
	"time"

	"github.com/SURF-Innovatie/MORIS/ent/event"
	"github.com/SURF-Innovatie/MORIS/internal/domain/commands"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/SURF-Innovatie/MORIS/internal/domain/projection"
	"github.com/google/uuid"
)

type MemoryStream struct {
	events []events.Event
}

func NewStream() *MemoryStream {
	return &MemoryStream{
		events: make([]events.Event, 0),
	}
}

func (s *MemoryStream) Append(id uuid.UUID, ev events.Event) {
	s.events = append(s.events, ev)
}

func (s *MemoryStream) Events(id uuid.UUID) []events.Event {
	// Simple filter by AggregateID
	var res []events.Event
	for _, e := range s.events {
		if e.AggregateID() == id {
			res = append(res, e)
		}
	}
	return res
}

func (s *MemoryStream) Reduce(id uuid.UUID) *entities.Project {
	return projection.Reduce(id, s.Events(id))
}

func Test_ProjectLifecycle(t *testing.T) {
	// Arrange
	stream := NewStream()
	id := uuid.New()

	start := time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 2, 2, 0, 0, 0, 0, time.UTC)
	org := uuid.New()
	actor := uuid.New()

	// Initial member
	personID := uuid.New()
	roleID := uuid.New()
	members := []entities.ProjectMember{
		{PersonID: personID, ProjectRoleID: roleID},
	}

	// StartProject
	ev, err := commands.StartProject(id, actor, "Alpha", "First", start, end, members, org)
	if err != nil {
		t.Fatal(err)
	}
	if ev == nil {
		t.Fatal("expected start event")
	}
	stream.Append(id, ev)

	// Reduce -> check initial state
	cur := stream.Reduce(id)
	if cur.Title != "Alpha" {
		t.Fatalf("title got %q", cur.Title)
	}
	if cur.Description != "First" {
		t.Fatalf("desc got %q", cur.Description)
	}
	// Org name is not in projection unless hydrated from DB, but ID is there
	if cur.OwningOrgNodeID != org {
		t.Fatalf("org ID mismatch")
	}
	if len(cur.Members) != 1 || cur.Members[0].PersonID != personID {
		t.Fatalf("members mismatch")
	}

	// ChangeTitle
	ev, err = commands.ChangeTitle(id, actor, cur, "Alpha v2", event.StatusApproved)
	if err != nil {
		t.Fatal(err)
	}
	if ev == nil {
		t.Fatal("expected title change event")
	}
	stream.Append(id, ev)
	// Apply only last event to current projection
	projection.Apply(cur, ev)
	if cur.Title != "Alpha v2" {
		t.Fatalf("title got %q", cur.Title)
	}

	// No-op ChangeTitle
	ev, err = commands.ChangeTitle(id, actor, cur, "Alpha v2", event.StatusApproved)
	if err != nil {
		t.Fatal(err)
	}
	if ev != nil {
		t.Fatalf("expected no event on same title")
	}

	// ChangeStartDate
	newStart := start.AddDate(0, 0, 1)
	ev, err = commands.ChangeStartDate(id, actor, cur, newStart, event.StatusApproved)
	if err != nil {
		t.Fatal(err)
	}
	if ev == nil {
		t.Fatal("expected start date changed")
	}
	stream.Append(id, ev)
	projection.Apply(cur, ev)
	if !cur.StartDate.Equal(newStart) {
		t.Fatalf("start date not updated")
	}

	// AddPerson (AssignProjectRole)
	newPerson := uuid.New()
	newRole := uuid.New()
	ev, err = commands.AssignProjectRole(id, actor, cur, newPerson, newRole, event.StatusApproved)
	if err != nil {
		t.Fatal(err)
	}
	if ev == nil {
		t.Fatal("expected role assigned")
	}
	stream.Append(id, ev)
	projection.Apply(cur, ev)

	if !hasMember(cur, newPerson) {
		t.Fatalf("newPerson not present")
	}

	// RemovePerson (UnassignProjectRole)
	ev, err = commands.UnassignProjectRole(id, actor, cur, newPerson, newRole, event.StatusApproved)
	if err != nil {
		t.Fatal(err)
	}
	if ev == nil {
		t.Fatal("expected role unassigned")
	}
	stream.Append(id, ev)
	projection.Apply(cur, ev)

	if hasMember(cur, newPerson) {
		t.Fatalf("newPerson still present")
	}

	// Re-reduce from full stream should equal current projection
	full := projection.Reduce(id, stream.Events(id))
	if full.Title != cur.Title ||
		!full.StartDate.Equal(cur.StartDate) ||
		len(full.Members) != len(cur.Members) {
		t.Fatalf("full reduce diverged from incremental apply")
	}
}

func hasMember(p *entities.Project, personID uuid.UUID) bool {
	for _, m := range p.Members {
		if m.PersonID == personID {
			return true
		}
	}
	return false
}
