package events

import (
	"github.com/google/uuid"
)

type OrganisationChanged struct {
	Base
	OrganisationID uuid.UUID `json:"organisation"`
}

func (OrganisationChanged) isEvent()     {}
func (OrganisationChanged) Type() string { return OrganisationChangedType }
