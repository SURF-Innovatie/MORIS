package events

import (
	"fmt"

	"github.com/google/uuid"
)

type PersonAdded struct {
	Base
	PersonId uuid.UUID `json:"personId"`
}

func (PersonAdded) isEvent()     {}
func (PersonAdded) Type() string { return PersonAddedType }
func (e PersonAdded) String() string {
	return fmt.Sprintf("Person added: %s", e.PersonId)
}
