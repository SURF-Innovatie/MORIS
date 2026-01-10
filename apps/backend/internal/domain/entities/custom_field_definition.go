package entities

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/customfielddefinition"
	"github.com/google/uuid"
)

type CustomFieldDefinition struct {
	ID                 uuid.UUID                      `json:"id"`
	OrganisationNodeID uuid.UUID                      `json:"organisation_node_id"`
	Name               string                         `json:"name"`
	Type               customfielddefinition.Type     `json:"type"`
	Category           customfielddefinition.Category `json:"category"`
	Description        string                         `json:"description"`
	Required           bool                           `json:"required"`
	ValidationRegex    string                         `json:"validation_regex"`
	ExampleValue       string                         `json:"example_value"`
}

func (p *CustomFieldDefinition) FromEnt(row *ent.CustomFieldDefinition) *CustomFieldDefinition {
	return &CustomFieldDefinition{
		ID:                 row.ID,
		OrganisationNodeID: row.OrganisationNodeID,
		Name:               row.Name,
		Type:               row.Type,
		Category:           row.Category,
		Description:        row.Description,
		Required:           row.Required,
		ValidationRegex:    row.ValidationRegex,
		ExampleValue:       row.ExampleValue,
	}
}
