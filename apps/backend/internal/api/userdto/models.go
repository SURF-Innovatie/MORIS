package userdto

import (
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type Request struct {
	PersonID uuid.UUID `json:"person_id"`
	Password string    `json:"password"`
}

type Response struct {
	ID          uuid.UUID `json:"id"`
	PersonID    uuid.UUID `json:"person_id"`
	ORCiD       *string   `json:"orcid"`
	Name        string    `json:"name"`
	GivenName   *string   `json:"givenName"`
	FamilyName  *string   `json:"familyName"`
	Email       string    `json:"email"`
	AvatarURL   *string   `json:"avatarUrl"`
	Description *string   `json:"description"`
}

func FromEntity(acc *entities.UserAccount) Response {
	p := acc.Person
	u := acc.User

	return Response{
		ID:          u.ID,
		PersonID:    u.PersonID,
		Email:       p.Email,
		Name:        p.Name,
		ORCiD:       p.ORCiD,
		GivenName:   p.GivenName,
		FamilyName:  p.FamilyName,
		AvatarURL:   p.AvatarUrl,
		Description: p.Description,
	}
}
