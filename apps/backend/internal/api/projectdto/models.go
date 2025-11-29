package projectdto

import (
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/api/organisationdto"
	"github.com/SURF-Innovatie/MORIS/internal/api/persondto"
	"github.com/SURF-Innovatie/MORIS/internal/api/productdto"
	"github.com/google/uuid"
)

// Request represents the request body for starting a new project
type Request struct {
	ProjectAdmin   uuid.UUID `json:"projectAdmin"`
	Title          string    `json:"title" example:"New Project"`
	Description    string    `json:"description" example:"This is a new project"`
	OrganisationID uuid.UUID `json:"organisationID"`
	StartDate      string    `json:"startDate" example:"2025-01-01T00:00:00Z"`
	EndDate        string    `json:"endDate" example:"2025-12-31T23:59:59Z"`
}

type Response struct {
	Id           uuid.UUID                `json:"id"`
	ProjectAdmin uuid.UUID                `json:"projectAdmin"`
	Version      int                      `json:"version"`
	Title        string                   `json:"title" example:"New Project"`
	Description  string                   `json:"description" example:"This is a new project"`
	StartDate    time.Time                `json:"startDate" example:"2025-01-01T00:00:00Z"`
	EndDate      time.Time                `json:"endDate" example:"2025-12-31T23:59:59Z"`
	Organization organisationdto.Response `json:"organization"`
	People       []persondto.Response     `json:"people"`
	Products     []productdto.Response    `json:"products"`
}
