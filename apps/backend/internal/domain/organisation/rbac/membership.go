package rbac

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/google/uuid"
)

type Membership struct {
	ID          uuid.UUID
	PersonID    uuid.UUID
	RoleScopeID uuid.UUID
}

func (p *Membership) FromEnt(row *ent.Membership) *Membership {
	return &Membership{
		ID:          row.ID,
		PersonID:    row.PersonID,
		RoleScopeID: row.RoleScopeID,
	}
}
