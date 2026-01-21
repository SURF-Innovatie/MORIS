package entities

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/google/uuid"
)

type Portfolio struct {
	ID               uuid.UUID
	PersonID         uuid.UUID
	Headline         *string
	Summary          *string
	Website          *string
	ShowEmail        bool
	ShowOrcid        bool
	PinnedProjectIDs []uuid.UUID
	PinnedProductIDs []uuid.UUID
}

func (p *Portfolio) FromEnt(row *ent.Portfolio) *Portfolio {
	return &Portfolio{
		ID:               row.ID,
		PersonID:         row.PersonID,
		Headline:         row.Headline,
		Summary:          row.Summary,
		Website:          row.Website,
		ShowEmail:        row.ShowEmail,
		ShowOrcid:        row.ShowOrcid,
		PinnedProjectIDs: row.PinnedProjectIds,
		PinnedProductIDs: row.PinnedProductIds,
	}
}
