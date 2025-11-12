package events

import "github.com/SURF-Innovatie/MORIS/internal/domain/entities"

type OrganisationChanged struct {
	Base
	Organisation entities.Organisation `json:"organisation"`
}

func (OrganisationChanged) isEvent()     {}
func (OrganisationChanged) Type() string { return OrganisationChangedType }
