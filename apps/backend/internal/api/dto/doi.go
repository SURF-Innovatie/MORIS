package dto

import "github.com/SURF-Innovatie/MORIS/internal/domain/entities"

// Work represents the simplified product data resolved from a DOI
type Work struct {
	DOI       string               `json:"doi"`
	Title     string               `json:"title"`
	Type      entities.ProductType `json:"type"`
	Date      string               `json:"date,omitempty"`
	Publisher string               `json:"publisher,omitempty"`
	Authors   []string             `json:"authors,omitempty"`
}
