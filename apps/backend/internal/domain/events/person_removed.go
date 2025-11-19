package events

import "github.com/SURF-Innovatie/MORIS/internal/domain/entities"

type PersonRemoved struct {
	Base
	Person entities.Person `json:"person"`
}

func (PersonRemoved) isEvent()     {}
func (PersonRemoved) Type() string { return PersonRemovedType }
