package dto

import (
	"github.com/SURF-Innovatie/MORIS/external/doi"
	"github.com/SURF-Innovatie/MORIS/internal/domain/product"
)

type WorkAuthor struct {
	Given  string `json:"given,omitempty"`
	Family string `json:"family,omitempty"`
	Name   string `json:"name,omitempty"`
	ORCID  string `json:"orcid,omitempty"`
}

type Work struct {
	DOI       doi.DOI             `json:"doi"`
	Title     string              `json:"title"`
	Type      product.ProductType `json:"type"`
	Date      string              `json:"date,omitempty"`
	Publisher string              `json:"publisher,omitempty"`
	Authors   []WorkAuthor        `json:"authors,omitempty"`
}
