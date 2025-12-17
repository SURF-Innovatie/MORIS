package dto

import (
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type OrganisationCreateRootRequest struct {
	Name string `json:"name"`
}

type OrganisationCreateChildRequest struct {
	Name string `json:"name"`
}

type OrganisationUpdateRequest struct {
	Name     string     `json:"name"`
	ParentID *uuid.UUID `json:"parentId"` // null => root
}

type OrganisationResponse struct {
	ID       uuid.UUID  `json:"id"`
	ParentID *uuid.UUID `json:"parentId"`
	Name     string     `json:"name"`
}

func (r OrganisationResponse) FromEntity(n entities.OrganisationNode) OrganisationResponse {
	return OrganisationResponse{
		ID:       n.ID,
		ParentID: n.ParentID,
		Name:     n.Name,
	}
}
