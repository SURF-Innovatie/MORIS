package userdto

import "github.com/google/uuid"

type Request struct {
	PersonID uuid.UUID `json:"person_id"`
	Password string    `json:"password"`
}

type Response struct {
	ID         uuid.UUID `json:"id"`
	PersonID   uuid.UUID `json:"person_id"`
	ORCiD      *string   `json:"orcid"`
	Name       string    `json:"name"`
	GivenName  *string   `json:"givenName"`
	FamilyName *string   `json:"familyName"`
	Email      string    `json:"email"`
}
