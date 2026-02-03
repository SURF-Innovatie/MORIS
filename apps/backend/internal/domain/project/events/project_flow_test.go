package events_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/domain/project"
	events2 "github.com/SURF-Innovatie/MORIS/internal/domain/project/events"
	"github.com/SURF-Innovatie/MORIS/internal/domain/project/projection"
	"github.com/google/uuid"
)

type MemoryStream struct {
	events []events2.Event
}

func NewStream() *MemoryStream {
	return &MemoryStream{events: make([]events2.Event, 0)}
}

func (s *MemoryStream) Append(ev events2.Event) {
	s.events = append(s.events, ev)
}

func (s *MemoryStream) Events(id uuid.UUID) []events2.Event {
	var res []events2.Event
	for _, e := range s.events {
		if e.AggregateID() == id {
			res = append(res, e)
		}
	}
	return res
}

func (s *MemoryStream) Reduce(id uuid.UUID) *project.Project {
	return projection.Reduce(id, s.Events(id))
}

func Test_ProjectLifecycle(t *testing.T) {
	stream := NewStream()
	id := uuid.New()

	start := time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 2, 2, 0, 0, 0, 0, time.UTC)
	org := uuid.New()
	actor := uuid.New()

	// Initial member
	personID := uuid.New()
	roleID := uuid.New()
	members := []project.Member{
		{PersonID: personID, ProjectRoleID: roleID},
	}

	var ev events2.Event
	var err error

	// StartProject (decider)
	ev, err = events2.DecideProjectStarted(
		id,
		actor,
		events2.ProjectStartedInput{
			Title:           "Alpha",
			Description:     "First",
			StartDate:       start,
			EndDate:         end,
			Members:         members,
			OwningOrgNodeID: org,
		},
		events2.StatusApproved,
	)
	if err != nil {
		t.Fatal(err)
	}
	if ev == nil {
		t.Fatal("expected start event")
	}
	stream.Append(ev)

	// Reduce -> check initial state
	cur := stream.Reduce(id)
	if cur.Title != "Alpha" {
		t.Fatalf("title got %q", cur.Title)
	}
	if cur.Description != "First" {
		t.Fatalf("desc got %q", cur.Description)
	}
	if cur.OwningOrgNodeID != org {
		t.Fatalf("org ID mismatch")
	}
	if len(cur.Members) != 1 || cur.Members[0].PersonID != personID {
		t.Fatalf("members mismatch")
	}

	// ChangeTitle (decider)
	ev, err = events2.DecideTitleChanged(
		id,
		actor,
		cur,
		events2.TitleChangedInput{Title: "Alpha v2"},
		events2.StatusApproved,
	)
	if err != nil {
		t.Fatal(err)
	}
	if ev == nil {
		t.Fatal("expected title change event")
	}
	stream.Append(ev)
	projection.Apply(cur, ev)
	if cur.Title != "Alpha v2" {
		t.Fatalf("title got %q", cur.Title)
	}

	// No-op ChangeTitle
	ev, err = events2.DecideTitleChanged(
		id,
		actor,
		cur,
		events2.TitleChangedInput{Title: "Alpha v2"},
		events2.StatusApproved,
	)
	fmt.Printf("e == nil? %v, type=%T\n", ev == nil, ev)
	if err != nil {
		t.Fatal(err)
	}
	if ev != nil {
		t.Fatalf("expected no event on same title")
	}

	// ChangeStartDate
	newStart := start.AddDate(0, 0, 1)
	ev, err = events2.DecideStartDateChanged(
		id,
		actor,
		cur,
		events2.StartDateChangedInput{StartDate: newStart},
		events2.StatusApproved,
	)
	if err != nil {
		t.Fatal(err)
	}
	if ev == nil {
		t.Fatal("expected start date changed")
	}
	stream.Append(ev)
	projection.Apply(cur, ev)
	if !cur.StartDate.Equal(newStart) {
		t.Fatalf("start date not updated")
	}

	// AddPerson (AssignProjectRole)
	newPerson := uuid.New()
	newRole := uuid.New()

	ev, err = events2.DecideProjectRoleAssigned(
		id,
		actor,
		cur,
		events2.ProjectRoleAssignedInput{
			PersonID:      newPerson,
			ProjectRoleID: newRole,
		},
		events2.StatusApproved,
	)
	if err != nil {
		t.Fatal(err)
	}
	if ev == nil {
		t.Fatal("expected role assigned")
	}
	stream.Append(ev)
	projection.Apply(cur, ev)

	if !hasMember(cur, newPerson) {
		t.Fatalf("newPerson not present")
	}

	// RemovePerson (UnassignProjectRole)
	ev, err = events2.DecideProjectRoleUnassigned(
		id,
		actor,
		cur,
		events2.ProjectRoleUnassignedInput{
			PersonID:      newPerson,
			ProjectRoleID: newRole,
		},
		events2.StatusApproved,
	)
	if err != nil {
		t.Fatal(err)
	}
	if ev == nil {
		t.Fatal("expected role unassigned")
	}
	stream.Append(ev)
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

func hasMember(p *project.Project, personID uuid.UUID) bool {
	for _, m := range p.Members {
		if m.PersonID == personID {
			return true
		}
	}
	return false
}
