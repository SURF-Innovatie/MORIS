package entities

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/google/uuid"
)

type Person struct {
	ID          uuid.UUID
	AvatarUrl       *string
	Description     *string
	UserID          uuid.UUID
	ORCiD           *string
	Name            string
	GivenName       *string
	FamilyName      *string
	Email           string
	OrgCustomFields map[string]interface{} `json:"org_custom_fields"`
}

func (p *Person) FromEnt(row *ent.Person) *Person {
	return &Person{
		ID:              row.ID,
		Name:            row.Name,
		ORCiD:           &row.OrcidID,
		GivenName:       row.GivenName,
		FamilyName:      row.FamilyName,
		Email:           row.Email,
		AvatarUrl:       row.AvatarURL,
		Description:     row.Description,
		OrgCustomFields: row.OrgCustomFields,
	}
}
