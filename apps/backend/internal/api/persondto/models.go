package persondto

import (
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type Request struct {
	UserID      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	GivenName   *string   `json:"givenName"`
	FamilyName  *string   `json:"familyName"`
	Email       string    `json:"email"`
	ORCiD       *string   `json:"orcid"`
	AvatarUrl   *string   `json:"avatar_url"`
	Description *string   `json:"description"`
}

type Response struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	GivenName   *string   `json:"givenName"`
	FamilyName  *string   `json:"familyName"`
	Email       string    `json:"email"`
	ORCiD       *string   `json:"orcid"`
	AvatarUrl   *string   `json:"avatar_url"`
	Description *string   `json:"description"`
}

func FromEntity(e entities.Person) Response {
	return Response{
		ID:          e.Id,
		UserID:      e.UserID,
		Name:        e.Name,
		GivenName:   e.GivenName,
		FamilyName:  e.FamilyName,
		Email:       e.Email,
		ORCiD:       e.ORCiD,
		AvatarUrl:   e.AvatarUrl,
		Description: e.Description,
	}
}
