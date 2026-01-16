package events_test

import (
	"testing"

	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
)

func TestAllRegisteredEventsHaveDeciders(t *testing.T) {
	if err := events.ValidateRegistrations(); err != nil {
		t.Fatal(err)
	}
}
