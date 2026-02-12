package queries

import (
	"github.com/SURF-Innovatie/MORIS/internal/domain/organisation"
	"github.com/SURF-Innovatie/MORIS/internal/domain/product"
	"github.com/SURF-Innovatie/MORIS/internal/domain/project"
)

type ProjectDetails struct {
	Project       project.Project               `json:"project"`
	OwningOrgNode organisation.OrganisationNode `json:"owning_org_node"`
	Members       []project.MemberDetail        `json:"members"`
	Products      []product.Product             `json:"products"`
}
