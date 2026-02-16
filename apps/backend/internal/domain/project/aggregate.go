package project

import (
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/domain/identity"
	"github.com/SURF-Innovatie/MORIS/internal/domain/project/role"
	"github.com/google/uuid"
)

type Member struct {
	PersonID      uuid.UUID `json:"person_id"`
	ProjectRoleID uuid.UUID `json:"project_role_id"`
}

type Project struct {
	Id                        uuid.UUID              `json:"id"`
	Version                   int                    `json:"version"`
	StartDate                 time.Time              `json:"start_date"`
	EndDate                   time.Time              `json:"end_date"`
	Title                     string                 `json:"title"`
	Description               string                 `json:"description"`
	Members                   []Member               `json:"members"`
	OwningOrgNodeID           uuid.UUID              `json:"owning_org_node_id"`
	ProductIDs                []uuid.UUID            `json:"product_ids"`
	AffiliatedOrganisationIDs []uuid.UUID            `json:"affiliated_organisation_ids"`
	CustomFields              map[string]interface{} `json:"custom_fields"`
}

type MemberDetail struct {
	Person identity.Person
	Role   role.ProjectRole
}

type ChangeLogEntry struct {
	Event string
	At    time.Time
}

type ChangeLog struct {
	Entries []ChangeLogEntry
}
