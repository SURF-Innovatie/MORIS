package events

import (
	"fmt"

	"github.com/google/uuid"
)

type PersonRemoved struct {
	Base
	PersonId uuid.UUID `json:"personId"`
}

func (PersonRemoved) isEvent()     {}
func (PersonRemoved) Type() string { return PersonRemovedType }
func (e PersonRemoved) String() string {
	// TODO: have person removed event contain edge to person
	return fmt.Sprintf("Person removed: %s", e.PersonId)
}
