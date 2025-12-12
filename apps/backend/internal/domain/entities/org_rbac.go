package entities

import "github.com/google/uuid"

type OrganisationRole struct {
	ID             uuid.UUID
	Key            string
	HasAdminRights bool
}

type RoleScope struct {
	ID         uuid.UUID
	RoleID     uuid.UUID
	RootNodeID uuid.UUID
}

type Membership struct {
	ID          uuid.UUID
	PersonID    uuid.UUID
	RoleScopeID uuid.UUID
}
