package entities

import "github.com/google/uuid"

type Person struct {
	Id   uuid.UUID
	Name string
}

// NewPerson creates a new person with a UUID
func NewPerson(name string) *Person {
	return &Person{
		Id:   uuid.New(),
		Name: name,
	}
}
