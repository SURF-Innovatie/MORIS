package entities

import "github.com/google/uuid"

type Person struct {
	Id   uuid.UUID
	Name string
}
