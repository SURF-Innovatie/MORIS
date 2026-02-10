package bulkimport

import (
	"github.com/SURF-Innovatie/MORIS/internal/api/dto"
	"github.com/google/uuid"
)

type Entry struct {
	DOI string `json:"doi"`
}

type ItemResult struct {
	DOI       string    `json:"doi"`
	Work      *dto.Work `json:"work,omitempty"`
	ProductID uuid.UUID `json:"product_id,omitempty"`
	Error     string    `json:"error,omitempty"`
}

type Result struct {
	ProjectID       uuid.UUID    `json:"project_id"`
	CreatedProducts []uuid.UUID  `json:"created_products"`
	Errors          []EntryError `json:"errors"`
	Items           []ItemResult `json:"items"`
}

type EntryError struct {
	DOI   string `json:"doi"`
	Error string `json:"error"`
}
