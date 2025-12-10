package eventdto

import (
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/api/persondto"
	"github.com/SURF-Innovatie/MORIS/internal/api/productdto"
	"github.com/google/uuid"
)

type Event struct {
	ID           uuid.UUID            `json:"id"`
	ProjectID    uuid.UUID            `json:"projectId"`
	Type         string               `json:"type"`
	Status       string               `json:"status"`
	CreatedBy    uuid.UUID            `json:"createdBy"`
	At           time.Time            `json:"at"`
	Details      string               `json:"details"`      // Human readable description
	ProjectTitle string               `json:"projectTitle"` // Title of the project
	Person       *persondto.Response  `json:"person,omitempty"`
	Product      *productdto.Response `json:"product,omitempty"`
}

type Response struct {
	Events []Event `json:"events"`
}
