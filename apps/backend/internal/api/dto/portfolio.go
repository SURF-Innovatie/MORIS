package dto

import (
	"github.com/SURF-Innovatie/MORIS/internal/domain/portfolio"
	"github.com/google/uuid"
)

type PortfolioRequest struct {
	Headline         *string     `json:"headline,omitempty"`
	Summary          *string     `json:"summary,omitempty"`
	Website          *string     `json:"website,omitempty"`
	ShowEmail        *bool       `json:"show_email,omitempty"`
	ShowOrcid        *bool       `json:"show_orcid,omitempty"`
	PinnedProjectIDs []uuid.UUID `json:"pinned_project_ids,omitempty"`
	PinnedProductIDs []uuid.UUID `json:"pinned_product_ids,omitempty"`
}

type PortfolioResponse struct {
	ID               uuid.UUID   `json:"id"`
	PersonID         uuid.UUID   `json:"person_id"`
	Headline         *string     `json:"headline,omitempty"`
	Summary          *string     `json:"summary,omitempty"`
	Website          *string     `json:"website,omitempty"`
	ShowEmail        bool        `json:"show_email"`
	ShowOrcid        bool        `json:"show_orcid"`
	PinnedProjectIDs []uuid.UUID `json:"pinned_project_ids"`
	PinnedProductIDs []uuid.UUID `json:"pinned_product_ids"`
	RecentProjectIDs []uuid.UUID `json:"recent_project_ids"`
}

func (r PortfolioResponse) FromEntity(e portfolio.Portfolio) PortfolioResponse {
	return PortfolioResponse{
		ID:               e.ID,
		PersonID:         e.PersonID,
		Headline:         e.Headline,
		Summary:          e.Summary,
		Website:          e.Website,
		ShowEmail:        e.ShowEmail,
		ShowOrcid:        e.ShowOrcid,
		PinnedProjectIDs: e.PinnedProjectIDs,
		PinnedProductIDs: e.PinnedProductIDs,
		RecentProjectIDs: e.RecentProjectIDs,
	}
}
