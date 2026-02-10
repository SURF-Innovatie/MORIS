package dto

import "github.com/google/uuid"

type BulkImportRequest struct {
	Dois []string `json:"dois"`
}

type BulkImportItemResult struct {
	DOI       string           `json:"doi"`
	ProductID *uuid.UUID       `json:"product_id,omitempty"`
	Product   *ProductResponse `json:"product,omitempty"`
	Error     *string          `json:"error,omitempty"`
}

type BulkImportResponse struct {
	ProjectID       uuid.UUID              `json:"project_id"`
	CreatedProducts []uuid.UUID            `json:"created_products"`
	Items           []BulkImportItemResult `json:"items"`
}
