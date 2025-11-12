package events

import "github.com/SURF-Innovatie/MORIS/internal/domain/entities"

type PersonAdded struct {
	Base
	Person entities.Person `json:"person"`
}

func (PersonAdded) isEvent()     {}
func (PersonAdded) Type() string { return PersonAddedType }
