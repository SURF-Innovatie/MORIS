package events

import (
	"github.com/google/uuid"
)

type PersonAdded struct {
	Base
	PersonId uuid.UUID `json:"personId"`
}

func (PersonAdded) isEvent()     {}
func (PersonAdded) Type() string { return PersonAddedType }
