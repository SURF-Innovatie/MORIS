package dto

import (
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/app/project/queries"
	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
	"github.com/SURF-Innovatie/MORIS/internal/domain/project"
	"github.com/SURF-Innovatie/MORIS/internal/domain/project/role"
	"github.com/google/uuid"
)

// ProjectRequest represents the request body for starting a new project
type ProjectRequest struct {
	ProjectAdmin    uuid.UUID `json:"project_admin"`
	Title           string    `json:"title" example:"NewService Project"`
	Description     string    `json:"description" example:"This is a new project"`
	OwningOrgNodeID uuid.UUID `json:"owning_org_node_id"`
	StartDate       string    `json:"start_date" example:"2025-01-01T00:00:00Z"`
	EndDate         string    `json:"end_date" example:"2025-12-31T23:59:59Z"`
}

type ProjectUpdateMemberRequest struct {
	Role string `json:"role" example:"contributor"`
}

type ProjectRoleCreateRequest struct {
	Key  string `json:"key" example:"specialist"`
	Name string `json:"name" example:"Specialist"`
}

type ProjectRoleUpdateRequest struct {
	AllowedEventTypes []string `json:"allowedEventTypes"`
}

type ProjectRoleResponse struct {
	ID                uuid.UUID `json:"id" example:"b990c264-b3c1-425f-88a1-69f22ba4a7a5"`
	Key               string    `json:"key" example:"contributor"`
	Name              string    `json:"name" example:"Contributor"`
	AllowedEventTypes []string  `json:"allowedEventTypes,omitempty"`
}

func (r ProjectRoleResponse) FromEntity(e role.ProjectRole) ProjectRoleResponse {
	return ProjectRoleResponse{
		ID:                e.ID,
		Key:               e.Key,
		Name:              e.Name,
		AllowedEventTypes: e.AllowedEventTypes,
	}
}

type ProjectMemberResponse struct {
	PersonResponse
	RoleID   uuid.UUID `json:"role_id"`
	Role     string    `json:"role"`
	RoleName string    `json:"role_name"`
}

func (r ProjectMemberResponse) FromEntity(e project.MemberDetail) ProjectMemberResponse {
	return ProjectMemberResponse{
		PersonResponse: transform.ToDTOItem[PersonResponse](e.Person),
		RoleID:         e.Role.ID,
		Role:           e.Role.Key,
		RoleName:       e.Role.Name,
	}
}

type ProjectResponse struct {
	Id                      uuid.UUID                        `json:"id"`
	Version                 int                              `json:"version"`
	Title                   string                           `json:"title" example:"NewService Project"`
	Description             string                           `json:"description" example:"This is a new project"`
	StartDate               time.Time                        `json:"start_date" example:"2025-01-01T00:00:00Z"`
	EndDate                 time.Time                        `json:"end_date" example:"2025-12-31T23:59:59Z"`
	OwningOrgNode           OrganisationResponse             `json:"owning_org_node"`
	Members                 []ProjectMemberResponse          `json:"members"`
	Products                []ProductResponse                `json:"products"`
	AffiliatedOrganisations []AffiliatedOrganisationResponse `json:"affiliated_organisations"`
	CustomFields            map[string]interface{}           `json:"custom_fields"`
}

func (r ProjectResponse) FromEntity(d *queries.ProjectDetails) ProjectResponse {
	return ProjectResponse{
		Id:                      d.Project.Id,
		Version:                 d.Project.Version,
		Title:                   d.Project.Title,
		Description:             d.Project.Description,
		StartDate:               d.Project.StartDate,
		EndDate:                 d.Project.EndDate,
		OwningOrgNode:           transform.ToDTOItem[OrganisationResponse](d.OwningOrgNode),
		Members:                 transform.ToDTOs[ProjectMemberResponse](d.Members),
		Products:                transform.ToDTOs[ProductResponse](d.Products),
		AffiliatedOrganisations: transform.ToDTOs[AffiliatedOrganisationResponse](d.AffiliatedOrganisations),
		CustomFields:            d.Project.CustomFields,
	}
}
