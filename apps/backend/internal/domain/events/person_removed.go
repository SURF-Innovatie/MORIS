package events

import (
	"github.com/google/uuid"
)

type PersonRemoved struct {
	Base
	PersonId uuid.UUID `json:"personId"`
}

func (PersonRemoved) isEvent()     {}
func (PersonRemoved) Type() string { return PersonRemovedType }
