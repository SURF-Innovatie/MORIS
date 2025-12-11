package entities

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/google/uuid"
)

type Person struct {
	Id          uuid.UUID
	AvatarUrl   *string
	Description *string
	UserID      uuid.UUID
	ORCiD       *string
	Name        string
	GivenName   *string
	FamilyName  *string
	Email       string
}

func (p *Person) FromEnt(row *ent.Person) *Person {
	return &Person{
		Id:          row.ID,
		Name:        row.Name,
		ORCiD:       &row.OrcidID,
		GivenName:   row.GivenName,
		FamilyName:  row.FamilyName,
		Email:       row.Email,
		AvatarUrl:   row.AvatarURL,
		Description: row.Description,
	}
}
