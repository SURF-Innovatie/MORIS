package entities

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/google/uuid"
)

type OrganisationNodeClosure struct {
	AncestorID   uuid.UUID
	DescendantID uuid.UUID
	Depth        int
}

func (p *OrganisationNodeClosure) FromEnt(row *ent.OrganisationNodeClosure) *OrganisationNodeClosure {
	return &OrganisationNodeClosure{
		AncestorID:   row.AncestorID,
		DescendantID: row.DescendantID,
		Depth:        row.Depth,
	}
}
