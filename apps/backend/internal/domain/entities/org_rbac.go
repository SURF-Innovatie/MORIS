package entities

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/google/uuid"
)

type OrganisationRole struct {
	ID             uuid.UUID
	Key            string
	HasAdminRights bool
}

func (p *OrganisationRole) FromEnt(row *ent.OrganisationRole) *OrganisationRole {
	return &OrganisationRole{
		ID:             row.ID,
		Key:            row.Key,
		HasAdminRights: row.HasAdminRights,
	}
}

type RoleScope struct {
	ID         uuid.UUID
	RoleID     uuid.UUID
	RootNodeID uuid.UUID
}

func (p *RoleScope) FromEnt(row *ent.RoleScope) *RoleScope {
	return &RoleScope{
		ID:         row.ID,
		RoleID:     row.RoleID,
		RootNodeID: row.RootNodeID,
	}
}

type Membership struct {
	ID          uuid.UUID
	PersonID    uuid.UUID
	RoleScopeID uuid.UUID
}

func (p *Membership) FromEnt(row *ent.Membership) *Membership {
	return &Membership{
		ID:          row.ID,
		PersonID:    row.PersonID,
		RoleScopeID: row.RoleScopeID,
	}
}
