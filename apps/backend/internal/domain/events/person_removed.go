package events

import (
	"fmt"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
)

type PersonRemoved struct {
	Base
	Person entities.Person `json:"person"`
}

func (PersonRemoved) isEvent()     {}
func (PersonRemoved) Type() string { return PersonRemovedType }
func (e PersonRemoved) String() string {
	if e.Person.Name != "" {
		return fmt.Sprintf("Person removed: %s", e.Person.Name)
	}
	return fmt.Sprintf("Person removed: %s", e.Person.Id)
}
