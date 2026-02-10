package dto

import "github.com/google/uuid"

type BulkImportRequest struct {
	Dois []string `json:"dois"`
}

type BulkImportWork struct {
	DOI       string   `json:"doi"`
	Title     string   `json:"title"`
	Publisher string   `json:"publisher"`
	Type      int      `json:"type"`    // product.ProductType as int
	Authors   []string `json:"authors"` // if present in dto.Work
	Date      string   `json:"date"`    // if present in dto.Work
	// TODO: Funders, Awards, ROR when dto.Work supports it
}

type BulkImportItemResult struct {
	DOI       string          `json:"doi"`
	Work      *BulkImportWork `json:"work,omitempty"`
	ProductID *uuid.UUID      `json:"product_id,omitempty"`
	Error     *string         `json:"error,omitempty"`
}

type BulkImportResponse struct {
	ProjectID       uuid.UUID              `json:"project_id"`
	CreatedProducts []uuid.UUID            `json:"created_products"`
	Items           []BulkImportItemResult `json:"items"`
}
