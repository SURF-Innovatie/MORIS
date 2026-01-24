package events

import (
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
)

// DetailedEvent wraps a domain event and includes the full entities
// referenced by its ID fields.
type DetailedEvent struct {
	Event Event

	// Optional related entities
	Person                 *entities.Person
	Product                *entities.Product
	ProjectRole            *entities.ProjectRole
	OrgNode                *entities.OrganisationNode
	AffiliatedOrganisation *entities.AffiliatedOrganisation
	Creator                *entities.Person
}
