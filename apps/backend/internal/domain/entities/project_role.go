package entities

import "github.com/google/uuid"

type ProjectRole struct {
	ID   uuid.UUID
	Key  string
	Name string
}
