package dto

import (
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type ProductRequest struct {
	Name               string `json:"name"`
	Type               int    `json:"type"`
	Language           string `json:"language"`
	DOI                string `json:"doi"`
	ZenodoDepositionID *int   `json:"zenodo_deposition_id,omitempty"`
}

type ProductResponse struct {
	ID                 uuid.UUID            `json:"id"`
	Name               string               `json:"name"`
	Type               entities.ProductType `json:"type"`
	Language           string               `json:"language"`
	DOI                string               `json:"doi"`
	ZenodoDepositionID int                  `json:"zenodo_deposition_id,omitempty"`
}

func (r ProductResponse) FromEntity(e entities.Product) ProductResponse {
	return ProductResponse{
		ID:                 e.Id,
		Name:               e.Name,
		Type:               e.Type,
		Language:           e.Language,
		DOI:                e.DOI,
		ZenodoDepositionID: e.ZenodoDepositionID,
	}
}
