package events

import (
	"fmt"
	"time"
)

type EndDateChanged struct {
	Base
	EndDate time.Time `json:"endDate"`
}

func (EndDateChanged) isEvent()     {}
func (EndDateChanged) Type() string { return EndDateChangedType }
func (e EndDateChanged) String() string {
	return fmt.Sprintf("End date changed: %s", e.EndDate)
}
