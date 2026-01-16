package entities

import (
	"slices"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/google/uuid"
)

type ProjectRole struct {
	ID                 uuid.UUID
	Key                string
	Name               string
	OrganisationNodeID uuid.UUID
	AllowedEventTypes  []string
}

func (p *ProjectRole) FromEnt(row *ent.ProjectRole) *ProjectRole {
	return &ProjectRole{
		ID:                 row.ID,
		Key:                row.Key,
		Name:               row.Name,
		OrganisationNodeID: row.OrganisationNodeID,
		AllowedEventTypes:  row.AllowedEventTypes,
	}
}

// CanUseEventType checks if this role is allowed to use the given event type
func (p *ProjectRole) CanUseEventType(eventType string) bool {
	if len(p.AllowedEventTypes) == 0 {
		return false // Empty means no events allowed
	}

	return slices.Contains(p.AllowedEventTypes, eventType)
}
