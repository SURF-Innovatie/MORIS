package dto

import (
	"github.com/SURF-Innovatie/MORIS/internal/domain/product"
)

// Work represents the simplified product data resolved from a DOI
type Work struct {
	DOI       string              `json:"doi"`
	Title     string              `json:"title"`
	Type      product.ProductType `json:"type"`
	Date      string              `json:"date,omitempty"`
	Publisher string              `json:"publisher,omitempty"`
	Authors   []string            `json:"authors,omitempty"`
}
