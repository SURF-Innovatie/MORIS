package raidsink_test

import (
	"testing"
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/domain/identity"
	"github.com/SURF-Innovatie/MORIS/internal/domain/organisation"
	events2 "github.com/SURF-Innovatie/MORIS/internal/domain/project/events"
	"github.com/google/uuid"

	"github.com/SURF-Innovatie/MORIS/internal/adapter"
	raidsink "github.com/SURF-Innovatie/MORIS/internal/adapter/sinks/raid"
)

func ptr(s string) *string { return &s }

func makeProjectStarted(projectID, actor uuid.UUID, title, desc string, start, end time.Time) *events2.ProjectStarted {
	evt := &events2.ProjectStarted{
		Title:       title,
		Description: desc,
		StartDate:   start,
		EndDate:     end,
	}
	evt.SetBase(events2.Base{
		ID:        uuid.New(),
		ProjectID: projectID,
		At:        start,
		CreatedBy: actor,
		Status:    events2.StatusApproved,
	})
	return evt
}

func makeTitleChanged(projectID, actor uuid.UUID, title string, at time.Time) *events2.TitleChanged {
	evt := &events2.TitleChanged{
		Title: title,
	}
	evt.SetBase(events2.Base{
		ID:        uuid.New(),
		ProjectID: projectID,
		At:        at,
		CreatedBy: actor,
		Status:    events2.StatusApproved,
	})
	return evt
}

func makeRoleAssigned(projectID, actor, personID, roleID uuid.UUID, at time.Time) *events2.ProjectRoleAssigned {
	evt := &events2.ProjectRoleAssigned{
		PersonID:      personID,
		ProjectRoleID: roleID,
	}
	evt.SetBase(events2.Base{
		ID:        uuid.New(),
		ProjectID: projectID,
		At:        at,
		CreatedBy: actor,
		Status:    events2.StatusApproved,
	})
	return evt
}

func makeRoleUnassigned(projectID, actor, personID, roleID uuid.UUID, at time.Time) *events2.ProjectRoleUnassigned {
	evt := &events2.ProjectRoleUnassigned{
		PersonID:      personID,
		ProjectRoleID: roleID,
	}
	evt.SetBase(events2.Base{
		ID:        uuid.New(),
		ProjectID: projectID,
		At:        at,
		CreatedBy: actor,
		Status:    events2.StatusApproved,
	})
	return evt
}

func TestRAiDMapper_TitleHistory(t *testing.T) {
	mapper := raidsink.NewRAiDMapper()

	projectID := uuid.New()
	actor := uuid.New()
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	titleChange := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

	pc := adapter.ProjectContext{
		ProjectID: projectID,
		Events: []events2.Event{
			makeProjectStarted(projectID, actor, "Original Title", "Description", start, end),
			makeTitleChanged(projectID, actor, "New Title", titleChange),
		},
	}

	req := mapper.MapToCreateRequest(pc)

	// Should have 2 titles in history
	if len(req.Title) != 2 {
		t.Fatalf("expected 2 titles (history), got %d", len(req.Title))
	}

	// First title should have end date
	if req.Title[0].Text != "Original Title" {
		t.Errorf("expected first title 'Original Title', got %q", req.Title[0].Text)
	}
	if req.Title[0].StartDate != "2024-01-01" {
		t.Errorf("expected first title start '2024-01-01', got %q", req.Title[0].StartDate)
	}
	if req.Title[0].EndDate == nil || *req.Title[0].EndDate != "2024-06-01" {
		t.Errorf("expected first title end '2024-06-01'")
	}

	// Second title should have no end date (current)
	if req.Title[1].Text != "New Title" {
		t.Errorf("expected second title 'New Title', got %q", req.Title[1].Text)
	}
	if req.Title[1].StartDate != "2024-06-01" {
		t.Errorf("expected second title start '2024-06-01', got %q", req.Title[1].StartDate)
	}
	if req.Title[1].EndDate != nil {
		t.Errorf("expected second title to have no end date (current)")
	}
}

func TestRAiDMapper_SingleTitle(t *testing.T) {
	mapper := raidsink.NewRAiDMapper()

	projectID := uuid.New()
	actor := uuid.New()
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

	pc := adapter.ProjectContext{
		ProjectID: projectID,
		Events: []events2.Event{
			makeProjectStarted(projectID, actor, "Test Project", "A test description", start, end),
		},
	}

	req := mapper.MapToCreateRequest(pc)

	// Should have 1 title
	if len(req.Title) != 1 {
		t.Fatalf("expected 1 title, got %d", len(req.Title))
	}
	if req.Title[0].Text != "Test Project" {
		t.Errorf("expected title 'Test Project', got %q", req.Title[0].Text)
	}
	// Current title should have no end date
	if req.Title[0].EndDate != nil {
		t.Errorf("expected no end date for current title")
	}
}

func TestRAiDMapper_ContributorsFromRoleEvents(t *testing.T) {
	mapper := raidsink.NewRAiDMapper()

	projectID := uuid.New()
	actor := uuid.New()
	personID := uuid.New()
	roleID := uuid.New()
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	assignedAt := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)

	pc := adapter.ProjectContext{
		ProjectID: projectID,
		Events: []events2.Event{
			makeProjectStarted(projectID, actor, "Test", "", start, start.AddDate(1, 0, 0)),
			makeRoleAssigned(projectID, actor, personID, roleID, assignedAt),
		},
		Members: []identity.Person{
			{
				ID:    personID,
				Name:  "John Doe",
				ORCiD: ptr("0000-0001-2345-6789"),
			},
		},
	}

	req := mapper.MapToCreateRequest(pc)

	// Should have 1 contributor
	if len(req.Contributor) != 1 {
		t.Fatalf("expected 1 contributor, got %d", len(req.Contributor))
	}
	if req.Contributor[0].Id == nil || *req.Contributor[0].Id != "0000-0001-2345-6789" {
		t.Errorf("expected ORCID '0000-0001-2345-6789'")
	}
	// Position should have start date from assignment event
	if len(req.Contributor[0].Position) != 1 {
		t.Fatalf("expected 1 position, got %d", len(req.Contributor[0].Position))
	}
	if req.Contributor[0].Position[0].StartDate != "2024-02-01" {
		t.Errorf("expected position start '2024-02-01', got %q", req.Contributor[0].Position[0].StartDate)
	}
}

func TestRAiDMapper_ContributorWithEndDate(t *testing.T) {
	mapper := raidsink.NewRAiDMapper()

	projectID := uuid.New()
	actor := uuid.New()
	personID := uuid.New()
	roleID := uuid.New()
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	assignedAt := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
	unassignedAt := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)

	pc := adapter.ProjectContext{
		ProjectID: projectID,
		Events: []events2.Event{
			makeProjectStarted(projectID, actor, "Test", "", start, start.AddDate(1, 0, 0)),
			makeRoleAssigned(projectID, actor, personID, roleID, assignedAt),
			makeRoleUnassigned(projectID, actor, personID, roleID, unassignedAt),
		},
		Members: []identity.Person{
			{
				ID:    personID,
				Name:  "John Doe",
				ORCiD: ptr("0000-0001-2345-6789"),
			},
		},
	}

	req := mapper.MapToCreateRequest(pc)

	if len(req.Contributor) != 1 {
		t.Fatalf("expected 1 contributor, got %d", len(req.Contributor))
	}

	// Should have end date from unassignment event
	pos := req.Contributor[0].Position[0]
	if pos.EndDate == nil {
		t.Fatal("expected end date for unassigned contributor")
	}
	if *pos.EndDate != "2024-06-01" {
		t.Errorf("expected end date '2024-06-01', got %q", *pos.EndDate)
	}
}

func TestRAiDMapper_SkipContributorWithoutORCID(t *testing.T) {
	mapper := raidsink.NewRAiDMapper()

	projectID := uuid.New()
	actor := uuid.New()
	personID := uuid.New()
	roleID := uuid.New()
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	pc := adapter.ProjectContext{
		ProjectID: projectID,
		Events: []events2.Event{
			makeProjectStarted(projectID, actor, "Test", "", start, start.AddDate(1, 0, 0)),
			makeRoleAssigned(projectID, actor, personID, roleID, start),
		},
		Members: []identity.Person{
			{
				ID:    personID,
				Name:  "Jane Smith",
				ORCiD: nil, // No ORCID
			},
		},
	}

	req := mapper.MapToCreateRequest(pc)

	// Should have 0 contributors (no ORCID)
	if len(req.Contributor) != 0 {
		t.Errorf("expected 0 contributors, got %d", len(req.Contributor))
	}
}

func TestRAiDMapper_OrganisationFromEvents(t *testing.T) {
	mapper := raidsink.NewRAiDMapper()

	projectID := uuid.New()
	actor := uuid.New()
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	pc := adapter.ProjectContext{
		ProjectID: projectID,
		Events: []events2.Event{
			makeProjectStarted(projectID, actor, "Test", "", start, start.AddDate(1, 0, 0)),
		},
		OrgNode: &organisation.OrganisationNode{
			ID:    uuid.New(),
			Name:  "Test University",
			RorID: ptr("https://ror.org/12345678"),
		},
	}

	req := mapper.MapToCreateRequest(pc)

	if len(req.Organisation) != 1 {
		t.Fatalf("expected 1 organisation, got %d", len(req.Organisation))
	}
	if req.Organisation[0].Id != "https://ror.org/12345678" {
		t.Errorf("expected ROR ID 'https://ror.org/12345678', got %q", req.Organisation[0].Id)
	}
}
