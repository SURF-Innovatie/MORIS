package dto

import (
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type OrganisationCreateRootRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	AvatarURL   *string `json:"avatarUrl"`
	// RorID is the Research Organization Registry ID
	RorID *string `json:"rorId"`
}

type OrganisationCreateChildRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	AvatarURL   *string `json:"avatarUrl"`
	// RorID is the Research Organization Registry ID
	RorID *string `json:"rorId"`
}

type OrganisationUpdateRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	AvatarURL   *string `json:"avatarUrl"`
	RorID       *string `json:"rorId"`
	// ParentID is the optional parent for moving the node
	ParentID *uuid.UUID `json:"parentId"` // null => root
}

type OrganisationResponse struct {
	ID          uuid.UUID  `json:"id"`
	ParentID    *uuid.UUID `json:"parentId"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	AvatarURL   *string    `json:"avatarUrl"`
	RorID       *string    `json:"rorId"`
}

func (r OrganisationResponse) FromEntity(n entities.OrganisationNode) OrganisationResponse {
	return OrganisationResponse{
		ID:          n.ID,
		ParentID:    n.ParentID,
		Name:        n.Name,
		Description: n.Description,
		AvatarURL:   n.AvatarURL,
		RorID:       n.RorID,
	}
}

type CustomFieldDefinitionCreateRequest struct {
	Name            string  `json:"name"`
	Type            string  `json:"type" example:"TEXT"` // TEXT, NUMBER, BOOLEAN, DATE
	Category        string  `json:"category" example:"PROJECT"` // PROJECT, PERSON
	Description     *string `json:"description"`
	Required        bool    `json:"required"`
	ValidationRegex *string `json:"validation_regex"`
	ExampleValue    *string `json:"example_value"`
}

type CustomFieldDefinitionResponse struct {
	ID                 uuid.UUID `json:"id"`
	OrganisationNodeID uuid.UUID `json:"organisation_node_id"`
	Name               string    `json:"name"`
	Type               string    `json:"type"`
	Category           string    `json:"category"`
	Description        string    `json:"description"`
	Required           bool      `json:"required"`
	ValidationRegex    string    `json:"validation_regex"`
	ExampleValue       string    `json:"example_value"`
}

type MemberCustomFieldUpdateValues struct {
	Values map[string]interface{} `json:"values"`
}
