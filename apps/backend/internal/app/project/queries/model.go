package queries

import (
	"github.com/SURF-Innovatie/MORIS/internal/domain/affiliatedorganisation"
	"github.com/SURF-Innovatie/MORIS/internal/domain/organisation"
	"github.com/SURF-Innovatie/MORIS/internal/domain/product"
	"github.com/SURF-Innovatie/MORIS/internal/domain/project"
)

type ProjectDetails struct {
	Project                 project.Project
	OwningOrgNode           organisation.OrganisationNode
	Members                 []project.MemberDetail
	Products                []product.Product
	AffiliatedOrganisations []affiliatedorganisation.AffiliatedOrganisation
}
