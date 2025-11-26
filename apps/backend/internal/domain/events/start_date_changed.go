package events

import (
	"fmt"
	"time"
)

type StartDateChanged struct {
	Base
	StartDate time.Time `json:"startDate"`
}

func (StartDateChanged) isEvent()     {}
func (StartDateChanged) Type() string { return StartDateChangedType }
func (e StartDateChanged) String() string {
	return fmt.Sprintf("Start date changed: %s", e.StartDate)
}
