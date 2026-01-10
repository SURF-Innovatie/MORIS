package entities

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/google/uuid"
)

type ProjectRole struct {
	ID                 uuid.UUID
	Key                string
	Name               string
	OrganisationNodeID uuid.UUID
}

func (p *ProjectRole) FromEnt(row *ent.ProjectRole) *ProjectRole {
	return &ProjectRole{
		ID:                 row.ID,
		Key:                row.Key,
		Name:               row.Name,
		OrganisationNodeID: row.OrganisationNodeID,
	}
}
