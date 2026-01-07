package dto

import (
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type OrganisationCreateRootRequest struct {
	Name string `json:"name"`
	// RorID is the Research Organization Registry ID
	RorID *string `json:"rorId"`
}

type OrganisationCreateChildRequest struct {
	Name string `json:"name"`
	// RorID is the Research Organization Registry ID
	RorID *string `json:"rorId"`
}

type OrganisationUpdateRequest struct {
	Name  string  `json:"name"`
	RorID *string `json:"rorId"`
	// ParentID is the optional parent for moving the node
	ParentID *uuid.UUID `json:"parentId"` // null => root
}

type OrganisationResponse struct {
	ID       uuid.UUID  `json:"id"`
	ParentID *uuid.UUID `json:"parentId"`
	Name     string     `json:"name"`
	RorID    *string    `json:"rorId"`
}

func (r OrganisationResponse) FromEntity(n entities.OrganisationNode) OrganisationResponse {
	return OrganisationResponse{
		ID:       n.ID,
		ParentID: n.ParentID,
		Name:     n.Name,
		RorID:    n.RorID,
	}
}
