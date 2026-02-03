package events

import (
	"github.com/SURF-Innovatie/MORIS/internal/domain/identity"
	"github.com/SURF-Innovatie/MORIS/internal/domain/organisation"
	"github.com/SURF-Innovatie/MORIS/internal/domain/product"
	"github.com/SURF-Innovatie/MORIS/internal/domain/project/role"
)

// DetailedEvent wraps a domain event and includes the full entities
// referenced by its ID fields.
type DetailedEvent struct {
	Event Event

	// Optional related entities
	Person      *identity.Person
	Product     *product.Product
	ProjectRole *role.ProjectRole
	OrgNode     *organisation.OrganisationNode
	Creator     *identity.Person
}
