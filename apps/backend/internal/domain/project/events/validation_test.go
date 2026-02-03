package events_test

import (
	"testing"
	"time"

	events2 "github.com/SURF-Innovatie/MORIS/internal/domain/project/events"
	"github.com/google/uuid"
)

func Test_StartProject_Validation(t *testing.T) {
	id := uuid.New()
	start := time.Now().UTC()
	end := start.Add(-time.Hour)
	actor := uuid.New()

	// missing title should error
	_, err := events2.DecideProjectStarted(
		id,
		actor,
		events2.ProjectStartedInput{
			Title:           "",
			Description:     "d",
			StartDate:       start,
			EndDate:         start,
			Members:         nil,
			OwningOrgNodeID: uuid.New(),
		},
		events2.StatusApproved,
	)
	if err == nil {
		t.Fatal("missing title should error")
	}

	// end before start should error
	_, err = events2.DecideProjectStarted(
		id,
		actor,
		events2.ProjectStartedInput{
			Title:           "t",
			Description:     "d",
			StartDate:       start,
			EndDate:         end,
			Members:         nil,
			OwningOrgNodeID: uuid.New(),
		},
		events2.StatusApproved,
	)
	if err == nil {
		t.Fatal("end before start should error")
	}
}
