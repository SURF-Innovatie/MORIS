package dto

import (
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type PersonRequest struct {
	UserID      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	GivenName   *string   `json:"givenName"`
	FamilyName  *string   `json:"familyName"`
	Email       string    `json:"email"`
	ORCiD       *string   `json:"orcid"`
	AvatarURL   *string   `json:"avatarUrl"`
	Description *string   `json:"description"`
}

type PersonResponse struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	GivenName   *string   `json:"givenName"`
	FamilyName  *string   `json:"familyName"`
	Email       string    `json:"email"`
	ORCiD       *string   `json:"orcid"`
	AvatarURL   *string   `json:"avatarUrl"`
	Description *string   `json:"description"`
}

func (r PersonResponse) FromEntity(e entities.Person) PersonResponse {
	return PersonResponse{
		ID:          e.ID,
		UserID:      e.UserID,
		Name:        e.Name,
		GivenName:   e.GivenName,
		FamilyName:  e.FamilyName,
		Email:       e.Email,
		ORCiD:       e.ORCiD,
		AvatarURL:   e.AvatarUrl,
		Description: e.Description,
	}
}
