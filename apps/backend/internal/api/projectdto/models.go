package projectdto

import (
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/api/organisationdto"
	"github.com/SURF-Innovatie/MORIS/internal/api/persondto"
	"github.com/SURF-Innovatie/MORIS/internal/api/productdto"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
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

func FromEntity(d entities.ProjectDetails) Response {
	peopleDTOs := make([]persondto.Response, 0, len(d.People))
	for _, p := range d.People {
		peopleDTOs = append(peopleDTOs, persondto.FromEntity(p))
	}

	productDTOs := make([]productdto.Response, 0, len(d.Products))
	for _, p := range d.Products {
		productDTOs = append(productDTOs, productdto.FromEntity(p))
	}

	return Response{
		Id:           d.Project.Id,
		ProjectAdmin: d.Project.ProjectAdmin,
		Version:      d.Project.Version,
		Title:        d.Project.Title,
		Description:  d.Project.Description,
		StartDate:    d.Project.StartDate,
		EndDate:      d.Project.EndDate,
		Organization: organisationdto.FromEntity(d.Organisation),
		People:       peopleDTOs,
		Products:     productDTOs,
	}
}
