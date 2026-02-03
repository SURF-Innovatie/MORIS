package rbac

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/google/uuid"
)

type OrganisationRole struct {
	ID                 uuid.UUID
	OrganisationNodeID uuid.UUID
	Key                string
	DisplayName        string
	Permissions        []Permission
}

func (p *OrganisationRole) FromEnt(row *ent.OrganisationRole) *OrganisationRole {
	perms := make([]Permission, len(row.Permissions))
	for i, v := range row.Permissions {
		perms[i] = Permission(v)
	}
	return &OrganisationRole{
		ID:                 row.ID,
		OrganisationNodeID: row.OrganisationNodeID,
		Key:                row.Key,
		DisplayName:        row.DisplayName,
		Permissions:        perms,
	}
}
