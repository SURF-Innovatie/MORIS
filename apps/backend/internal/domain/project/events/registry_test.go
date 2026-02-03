package events_test

import (
	"testing"

	"github.com/SURF-Innovatie/MORIS/internal/domain/project/events"
)

func TestAllRegisteredEventsHaveDeciders(t *testing.T) {
	if err := events.ValidateRegistrations(); err != nil {
		t.Fatal(err)
	}
}
