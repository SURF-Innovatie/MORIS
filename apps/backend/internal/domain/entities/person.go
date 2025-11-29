package entities

import "github.com/google/uuid"

type Person struct {
	Id         uuid.UUID
	UserID     uuid.UUID
	ORCiD      *string
	Name       string
	GivenName  *string
	FamilyName *string
	Email      string
}
