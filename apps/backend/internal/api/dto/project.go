package dto

import (
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

// ProjectRequest represents the request body for starting a new project
type ProjectRequest struct {
	ProjectAdmin    uuid.UUID `json:"project_admin"`
	Title           string    `json:"title" example:"New Project"`
	Description     string    `json:"description" example:"This is a new project"`
	OwningOrgNodeID uuid.UUID `json:"owning_org_node_id"`
	StartDate       string    `json:"start_date" example:"2025-01-01T00:00:00Z"`
	EndDate         string    `json:"end_date" example:"2025-12-31T23:59:59Z"`
}

type ProjectUpdateMemberRequest struct {
	Role string `json:"role" example:"contributor"`
}

type ProjectRoleResponse struct {
	Key  string `json:"key" example:"contributor"`
	Name string `json:"name" example:"Contributor"`
}

type ProjectMemberResponse struct {
	PersonResponse
	Roles []ProjectRoleResponse
}

func (r ProjectRoleResponse) FromEntity(e entities.ProjectRole) ProjectRoleResponse {
	return ProjectRoleResponse{
		Key:  e.Key,
		Name: e.Name,
	}
}

func (r ProjectMemberResponse) FromEntity(e entities.ProjectMemberDetail) ProjectMemberResponse {
	return ProjectMemberResponse{
		PersonResponse: transform.ToDTOItem[PersonResponse](e.Person),
		Roles:          transform.ToDTOs[ProjectRoleResponse](e.Roles),
	}
}

type ProjectResponse struct {
	Id            uuid.UUID               `json:"id"`
	Version       int                     `json:"version"`
	Title         string                  `json:"title" example:"New Project"`
	Description   string                  `json:"description" example:"This is a new project"`
	StartDate     time.Time               `json:"start_date" example:"2025-01-01T00:00:00Z"`
	EndDate       time.Time               `json:"end_date" example:"2025-12-31T23:59:59Z"`
	OwningOrgNode OrganisationResponse    `json:"owning_org_node"`
	Members       []ProjectMemberResponse `json:"members"`
	Products      []ProductResponse       `json:"products"`
}

func (r ProjectResponse) FromEntity(d *entities.ProjectDetails) ProjectResponse {
	return ProjectResponse{
		Id:            d.Project.Id,
		Version:       d.Project.Version,
		Title:         d.Project.Title,
		Description:   d.Project.Description,
		StartDate:     d.Project.StartDate,
		EndDate:       d.Project.EndDate,
		OwningOrgNode: transform.ToDTOItem[OrganisationResponse](d.OwningOrgNode),
		Members:       transform.ToDTOs[ProjectMemberResponse](d.Members),
		Products:      transform.ToDTOs[ProductResponse](d.Products),
	}
}
