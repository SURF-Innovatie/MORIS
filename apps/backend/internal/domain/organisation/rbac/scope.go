package rbac

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/google/uuid"
)

type RoleScope struct {
	ID         uuid.UUID
	RoleID     uuid.UUID
	RootNodeID uuid.UUID
}

func (p *RoleScope) FromEnt(row *ent.RoleScope) *RoleScope {
	return &RoleScope{
		ID:         row.ID,
		RoleID:     row.RoleID,
		RootNodeID: row.RootNodeID,
	}
}
