package persondto

import "github.com/google/uuid"

type Request struct {
	UserID     uuid.UUID `json:"user_id"`
	Name       string    `json:"name"`
	GivenName  *string   `json:"givenName"`
	FamilyName *string   `json:"familyName"`
	Email      string    `json:"email"`
	ORCiD      *string   `json:"orcid"`
}

type Response struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	Name       string    `json:"name"`
	GivenName  *string   `json:"givenName"`
	FamilyName *string   `json:"familyName"`
	Email      string    `json:"email"`
	ORCiD      *string   `json:"orcid"`
}
