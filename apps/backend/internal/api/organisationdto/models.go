package organisationdto

import (
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type CreateRootRequest struct {
	Name string `json:"name"`
}

type CreateChildRequest struct {
	Name string `json:"name"`
}

type UpdateRequest struct {
	Name     string     `json:"name"`
	ParentID *uuid.UUID `json:"parentId"` // null => root
}

type Response struct {
	ID       uuid.UUID  `json:"id"`
	ParentID *uuid.UUID `json:"parentId"`
	Name     string     `json:"name"`
}

func FromEntity(n entities.OrganisationNode) Response {
	return Response{
		ID:       n.ID,
		ParentID: n.ParentID,
		Name:     n.Name,
	}
}
