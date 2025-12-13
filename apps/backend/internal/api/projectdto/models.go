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
	ProjectAdmin    uuid.UUID `json:"project_admin"`
	Title           string    `json:"title" example:"New Project"`
	Description     string    `json:"description" example:"This is a new project"`
	OwningOrgNodeID uuid.UUID `json:"owning_org_node_id"`
	StartDate       string    `json:"start_date" example:"2025-01-01T00:00:00Z"`
	EndDate         string    `json:"end_date" example:"2025-12-31T23:59:59Z"`
}

type Response struct {
	Id            uuid.UUID                `json:"id"`
	Version       int                      `json:"version"`
	Title         string                   `json:"title" example:"New Project"`
	Description   string                   `json:"description" example:"This is a new project"`
	StartDate     time.Time                `json:"start_date" example:"2025-01-01T00:00:00Z"`
	EndDate       time.Time                `json:"end_date" example:"2025-12-31T23:59:59Z"`
	OwningOrgNode organisationdto.Response `json:"owning_org_node"`
	People        []persondto.Response     `json:"people"`
	Products      []productdto.Response    `json:"products"`
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
		Id:            d.Project.Id,
		Version:       d.Project.Version,
		Title:         d.Project.Title,
		Description:   d.Project.Description,
		StartDate:     d.Project.StartDate,
		EndDate:       d.Project.EndDate,
		OwningOrgNode: organisationdto.FromEntity(d.OwningOrgNode),
		People:        peopleDTOs,
		Products:      productDTOs,
	}
}
