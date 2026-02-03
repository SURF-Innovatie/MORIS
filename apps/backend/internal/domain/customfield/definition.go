package customfield

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/google/uuid"
)

type Definition struct {
	ID                 uuid.UUID `json:"id"`
	OrganisationNodeID uuid.UUID `json:"organisation_node_id"`
	Name               string    `json:"name"`
	Type               Type      `json:"type"`
	Category           Category  `json:"category"`
	Description        string    `json:"description"`
	Required           bool      `json:"required"`
	ValidationRegex    string    `json:"validation_regex"`
	ExampleValue       string    `json:"example_value"`
}

func (p *Definition) FromEnt(row *ent.CustomFieldDefinition) *Definition {
	return &Definition{
		ID:                 row.ID,
		OrganisationNodeID: row.OrganisationNodeID,
		Name:               row.Name,
		Type:               Type(row.Type),
		Category:           Category(row.Category),
		Description:        row.Description,
		Required:           row.Required,
		ValidationRegex:    row.ValidationRegex,
		ExampleValue:       row.ExampleValue,
	}
}
