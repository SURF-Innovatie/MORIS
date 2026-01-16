package entities

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/google/uuid"
)

type CustomFieldDefinition struct {
	ID                 uuid.UUID           `json:"id"`
	OrganisationNodeID uuid.UUID           `json:"organisation_node_id"`
	Name               string              `json:"name"`
	Type               CustomFieldType     `json:"type"`
	Category           CustomFieldCategory `json:"category"`
	Description        string              `json:"description"`
	Required           bool                `json:"required"`
	ValidationRegex    string              `json:"validation_regex"`
	ExampleValue       string              `json:"example_value"`
}

func (p *CustomFieldDefinition) FromEnt(row *ent.CustomFieldDefinition) *CustomFieldDefinition {
	return &CustomFieldDefinition{
		ID:                 row.ID,
		OrganisationNodeID: row.OrganisationNodeID,
		Name:               row.Name,
		Type:               CustomFieldType(row.Type),
		Category:           CustomFieldCategory(row.Category),
		Description:        row.Description,
		Required:           row.Required,
		ValidationRegex:    row.ValidationRegex,
		ExampleValue:       row.ExampleValue,
	}
}

type CustomFieldType string

const (
	CustomFieldTypeText    CustomFieldType = "TEXT"
	CustomFieldTypeNumber  CustomFieldType = "NUMBER"
	CustomFieldTypeBoolean CustomFieldType = "BOOLEAN"
	CustomFieldTypeDate    CustomFieldType = "DATE"
)

func (t CustomFieldType) Valid() bool {
	switch t {
	case CustomFieldTypeText, CustomFieldTypeNumber, CustomFieldTypeBoolean, CustomFieldTypeDate:
		return true
	default:
		return false
	}
}

type CustomFieldCategory string

const (
	CustomFieldCategoryProject CustomFieldCategory = "PROJECT"
	CustomFieldCategoryPerson  CustomFieldCategory = "PERSON"
)

func (c CustomFieldCategory) Valid() bool {
	switch c {
	case CustomFieldCategoryProject, CustomFieldCategoryPerson:
		return true
	default:
		return false
	}
}
