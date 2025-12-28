package events_test

import (
	"testing"
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/google/uuid"
)

func Test_StartProject_Validation(t *testing.T) {
	id := uuid.New()
	start := time.Now().UTC()
	end := start.Add(-time.Hour)
	actor := uuid.New()

	// missing title should error
	_, err := events.DecideProjectStarted(
		id,
		actor,
		events.ProjectStartedInput{
			Title:           "",
			Description:     "d",
			StartDate:       start,
			EndDate:         start,
			Members:         nil,
			OwningOrgNodeID: uuid.New(),
		},
		events.StatusApproved,
	)
	if err == nil {
		t.Fatal("missing title should error")
	}

	// end before start should error
	_, err = events.DecideProjectStarted(
		id,
		actor,
		events.ProjectStartedInput{
			Title:           "t",
			Description:     "d",
			StartDate:       start,
			EndDate:         end,
			Members:         nil,
			OwningOrgNodeID: uuid.New(),
		},
		events.StatusApproved,
	)
	if err == nil {
		t.Fatal("end before start should error")
	}
}
