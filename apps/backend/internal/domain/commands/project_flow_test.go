package commands_test

import (
	"testing"
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/platform/eventstore"
	"github.com/google/uuid"

	"github.com/SURF-Innovatie/MORIS/internal/domain/commands"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/projection"
)

func Test_ProjectLifecycle(t *testing.T) {
	// Arrange
	stream := eventstore.NewStream()
	id := uuid.New()

	start := time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 2, 2, 0, 0, 0, 0, time.UTC)
	org := uuid.New()
	people := []uuid.UUID{uuid.New()}
	actor := uuid.New()

	// StartProject
	ev, err := commands.StartProject(id, actor, "Alpha", "First", start, end, people, org)
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
	if cur.Organisation.Name != "Org" {
		t.Fatalf("org got %q", cur.Organisation.Name)
	}
	if len(cur.People) != 1 || cur.People[0].Name != "Ada" {
		t.Fatalf("people mismatch")
	}

	// ChangeTitle
	ev, err = commands.ChangeTitle(id, actor, cur, "Alpha v2")
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
	ev, err = commands.ChangeTitle(id, actor, cur, "Alpha v2")
	if err != nil {
		t.Fatal(err)
	}
	if ev != nil {
		t.Fatalf("expected no event on same title")
	}

	// ChangeStartDate
	newStart := start.AddDate(0, 0, 1)
	ev, err = commands.ChangeStartDate(id, actor, cur, newStart)
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

	// AddPerson + RemovePerson
	newPerson := uuid.New()
	ev, err = commands.AddPerson(id, actor, cur, newPerson)
	if err != nil {
		t.Fatal(err)
	}
	if ev == nil {
		t.Fatal("expected person added")
	}
	stream.Append(id, ev)
	projection.Apply(cur, ev)
	if !has(cur, "Grace") {
		t.Fatalf("Grace not present")
	}

	ev, err = commands.RemovePerson(id, actor, cur, newPerson)
	if err != nil {
		t.Fatal(err)
	}
	if ev == nil {
		t.Fatal("expected person removed")
	}
	stream.Append(id, ev)
	projection.Apply(cur, ev)
	if has(cur, "Grace") {
		t.Fatalf("Grace still present")
	}

	// Re-reduce from full stream should equal current projection
	full := projection.Reduce(id, stream.Events(id))
	if full.Title != cur.Title ||
		!full.StartDate.Equal(cur.StartDate) ||
		len(full.People) != len(cur.People) {
		t.Fatalf("full reduce diverged from incremental apply")
	}
}

func has(p *entities.Project, name string) bool {
	for _, x := range p.People {
		if x != nil && x.Name == name {
			return true
		}
	}
	return false
}
