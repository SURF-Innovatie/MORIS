package events

import "time"

type EndDateChanged struct {
	Base
	EndDate time.Time `json:"endDate"`
}

func (EndDateChanged) isEvent()     {}
func (EndDateChanged) Type() string { return EndDateChangedType }
