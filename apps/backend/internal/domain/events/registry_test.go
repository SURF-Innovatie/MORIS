package events

import "testing"

func TestAllRegisteredEventsHaveDeciders(t *testing.T) {
	if err := ValidateRegistrations(); err != nil {
		t.Fatal(err)
	}
}
