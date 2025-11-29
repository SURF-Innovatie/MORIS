package entities

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID
	PersonID uuid.UUID
	Password string
}
