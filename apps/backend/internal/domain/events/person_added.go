package events

import (
	"fmt"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
)

type PersonAdded struct {
	Base
	Person entities.Person `json:"person"`
}

func (PersonAdded) isEvent()     {}
func (PersonAdded) Type() string { return PersonAddedType }
func (e PersonAdded) String() string {
	if e.Person.Name != "" {
		return fmt.Sprintf("Person added: %s", e.Person.Name)
	}
	return fmt.Sprintf("Person added: %s", e.Person.Id)
}
