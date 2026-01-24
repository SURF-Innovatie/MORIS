package queries

import "github.com/SURF-Innovatie/MORIS/internal/domain/entities"

type ProjectDetails struct {
	Project                 entities.Project
	OwningOrgNode           entities.OrganisationNode
	Members                 []entities.ProjectMemberDetail
	Products                []entities.Product
	AffiliatedOrganisations []entities.AffiliatedOrganisation
}
