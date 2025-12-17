package entities

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/google/uuid"
)

type OrganisationNode struct {
	ID       uuid.UUID
	ParentID *uuid.UUID
	Name     string
}

func (o *OrganisationNode) FromEnt(row *ent.OrganisationNode) *OrganisationNode {
	if row == nil {
		return nil
	}

	var parentID *uuid.UUID
	if row.ParentID != nil {
		parentID = row.ParentID
	}

	return &OrganisationNode{
		ID:       row.ID,
		ParentID: parentID,
		Name:     row.Name,
	}
}
