package events

import "time"

type StartDateChanged struct {
	Base
	StartDate time.Time `json:"startDate"`
}

func (StartDateChanged) isEvent()     {}
func (StartDateChanged) Type() string { return StartDateChangedType }
