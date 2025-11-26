package events

import (
	"fmt"

	"github.com/google/uuid"
)

type OrganisationChanged struct {
	Base
	OrganisationID uuid.UUID `json:"organisation"`
}

func (OrganisationChanged) isEvent()     {}
func (OrganisationChanged) Type() string { return OrganisationChangedType }
func (e OrganisationChanged) String() string {
	return fmt.Sprintf("Organisation changed: %s", e.OrganisationID)
}
