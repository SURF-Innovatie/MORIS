package entities

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/google/uuid"
)

type Person struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	ORCiD           *string
	Name            string
	GivenName       *string
	FamilyName      *string
	Email           string
	AvatarUrl       *string
	Description     *string
	OrgCustomFields map[string]interface{} `json:"org_custom_fields"`
}

func (p *Person) FromEnt(row *ent.Person) *Person {
	out := &Person{
		ID:              row.ID,
		Name:            row.Name,
		GivenName:       row.GivenName,
		FamilyName:      row.FamilyName,
		Email:           row.Email,
		AvatarUrl:       row.AvatarURL,
		Description:     row.Description,
		ORCiD:           &row.OrcidID,
		OrgCustomFields: row.OrgCustomFields,
	}

	type userIDPtr interface{ GetUserID() *uuid.UUID }
	if v, ok := any(row).(userIDPtr); ok {
		if v.GetUserID() != nil {
			out.UserID = *v.GetUserID()
		} else {
			out.UserID = uuid.Nil
		}
	} else {
		out.UserID = row.UserID
	}

	return out
}
