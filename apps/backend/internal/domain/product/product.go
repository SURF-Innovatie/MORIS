package product

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/google/uuid"
)

// TODO: expand using: https://openaire-guidelines-for-cris-managers.readthedocs.io/en/v1.2.0/cerif_xml_product_entity.html

type Product struct {
	Id                 uuid.UUID
	Type               ProductType
	Language           string // TODO: Make language comply with spec IETF BCP 47, see: https://openaire-guidelines-for-cris-managers.readthedocs.io/en/v1.2.0/cerif_xml_product_entity.html
	Name               string
	DOI                string
	ZenodoDepositionID int
	AuthorPersonIDs    []uuid.UUID
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
	zenodoID := 0
	if row.ZenodoDepositionID != nil {
		zenodoID = *row.ZenodoDepositionID
	}
	lang := ""
	if row.Language != nil {
		lang = *row.Language
	}

	doi := ""
	if row.Doi != nil {
		doi = *row.Doi
	}

	authorIDs := make([]uuid.UUID, 0, len(row.Edges.Authors))
	for _, a := range row.Edges.Authors {
		if a != nil {
			authorIDs = append(authorIDs, a.ID)
		}
	}

	return &Product{
		Id:                 row.ID,
		Type:               ProductType(row.Type),
		Language:           lang,
		Name:               row.Name,
		DOI:                doi,
		ZenodoDepositionID: zenodoID,
		AuthorPersonIDs:    authorIDs,
	}
}
