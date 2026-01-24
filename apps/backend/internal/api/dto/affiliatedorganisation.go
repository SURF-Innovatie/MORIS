package dto

import (
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

// AffiliatedOrganisationRequest is the request body for creating/updating an affiliated organisation.
type AffiliatedOrganisationRequest struct {
	Name      string `json:"name"`
	KvkNumber string `json:"kvk_number,omitempty"`
	RorID     string `json:"ror_id,omitempty"`
	VatNumber string `json:"vat_number,omitempty"`
	City      string `json:"city,omitempty"`
	Country   string `json:"country,omitempty"`
}

// AffiliatedOrganisationResponse is the response for an affiliated organisation.
type AffiliatedOrganisationResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	KvkNumber string    `json:"kvk_number,omitempty"`
	RorID     string    `json:"ror_id,omitempty"`
	VatNumber string    `json:"vat_number,omitempty"`
	City      string    `json:"city,omitempty"`
	Country   string    `json:"country,omitempty"`
}

// FromEntity creates a response from a domain entity.
func (r AffiliatedOrganisationResponse) FromEntity(e entities.AffiliatedOrganisation) AffiliatedOrganisationResponse {
	return AffiliatedOrganisationResponse{
		ID:        e.ID,
		Name:      e.Name,
		KvkNumber: e.KvkNumber,
		RorID:     e.RorID,
		VatNumber: e.VatNumber,
		City:      e.City,
		Country:   e.Country,
	}
}
