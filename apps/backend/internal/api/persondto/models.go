package persondto

import "github.com/google/uuid"

type Request struct {
	Name       string  `json:"name"`
	GivenName  *string `json:"givenName"`
	FamilyName *string `json:"familyName"`
	Email      *string `json:"email"`
}

type Response struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	GivenName  *string   `json:"givenName"`
	FamilyName *string   `json:"familyName"`
	Email      *string   `json:"email"`
}
