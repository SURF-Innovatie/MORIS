package entities

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/google/uuid"
)

// AffiliatedOrganisation represents an organisation affiliated with a project.
type AffiliatedOrganisation struct {
	ID        uuid.UUID
	Name      string
	KvkNumber string
	RorID     string
	VatNumber string
	City      string
	Country   string
}

// FromEnt creates an AffiliatedOrganisation from an ent row.
func (a *AffiliatedOrganisation) FromEnt(row *ent.AffiliatedOrganisation) *AffiliatedOrganisation {
	return &AffiliatedOrganisation{
		ID:        row.ID,
		Name:      row.Name,
		KvkNumber: row.KvkNumber,
		RorID:     row.RorID,
		VatNumber: row.VatNumber,
		City:      row.City,
		Country:   row.Country,
	}
}
