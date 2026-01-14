package entities

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/google/uuid"
)

type User struct {
	ID                 uuid.UUID
	PersonID           uuid.UUID
	Password           string
	IsSysAdmin         bool `json:"is_sys_admin"`
	IsActive           bool `json:"is_active"`
	ZenodoAccessToken  *string
	ZenodoRefreshToken *string
}

func (u *User) FromEnt(row *ent.User) *User {
	return &User{
		ID:         row.ID,
		PersonID:   row.PersonID,
		IsSysAdmin: row.IsSysAdmin,
		IsActive:   row.IsActive,
	}
}
