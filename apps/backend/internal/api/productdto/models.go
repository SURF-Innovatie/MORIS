package productdto

import (
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type Request struct {
	Name     string `json:"name"`
	Type     int    `json:"type"`
	Language string `json:"language"`
	DOI      string `json:"doi"`
}

type Response struct {
	ID       uuid.UUID            `json:"id"`
	Name     string               `json:"name"`
	Type     entities.ProductType `json:"type"`
	Language string               `json:"language"`
	DOI      string               `json:"doi"`
}

func FromEntity(e entities.Product) Response {
	return Response{
		ID:       e.Id,
		Name:     e.Name,
		Type:     e.Type,
		Language: e.Language,
		DOI:      e.DOI,
	}
}
