package entities

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/google/uuid"
)

// TODO: expand using: https://openaire-guidelines-for-cris-managers.readthedocs.io/en/v1.2.0/cerif_xml_product_entity.html

type Product struct {
	Id       uuid.UUID
	Type     ProductType
	Language string // TODO: Make language comply with spec IETF BCP 47, see: https://openaire-guidelines-for-cris-managers.readthedocs.io/en/v1.2.0/cerif_xml_product_entity.html
	Name     string
	DOI      string
}

type ProductType int

// TODO: Make this more elaborate, see: https://openaire-guidelines-for-cris-managers.readthedocs.io/en/v1.2.0/cerif_xml_product_entity.html
const (
	CartographicMaterial ProductType = iota
	Dataset
	Image
	InteractiveResource
	LearningObject
	Other
	Software
	Sound
	Trademark
	Workflow
)

func (p *Product) FromEnt(row *ent.Product) *Product {
	return &Product{
		Id:       row.ID,
		Type:     ProductType(row.Type),
		Language: *row.Language,
		Name:     row.Name,
		DOI:      *row.Doi,
	}
}
