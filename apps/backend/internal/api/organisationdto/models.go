package organisationdto

import (
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type CreateRequest struct {
	Name string `json:"name"`
}

type Response struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func FromEntity(e entities.Organisation) Response {
	return Response{
		ID:   e.Id,
		Name: e.Name,
	}
}
