package bulkimport

import (
	"github.com/SURF-Innovatie/MORIS/internal/domain/product"
	"github.com/google/uuid"
)

type Entry struct {
	DOI string `json:"doi"`
}

type ItemResult struct {
	DOI     string           `json:"doi"`
	Product *product.Product `json:"product,omitempty"`
	Error   string           `json:"error"`
}

type Result struct {
	ProjectID       uuid.UUID    `json:"project_id"`
	CreatedProducts []uuid.UUID  `json:"created_products"`
	Items           []ItemResult `json:"items"`
	Errors          []EntryError `json:"errors"`
}

type EntryError struct {
	DOI   string `json:"doi"`
	Error string `json:"error"`
}
