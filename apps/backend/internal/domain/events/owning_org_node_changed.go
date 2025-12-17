package events

import (
	"github.com/google/uuid"
)

type OwningOrgNodeChanged struct {
	Base
	OwningOrgNodeID uuid.UUID `json:"owning_org_node_id"`
}

func (OwningOrgNodeChanged) isEvent()     {}
func (OwningOrgNodeChanged) Type() string { return OwningOrgNodeChangedType }
func (e OwningOrgNodeChanged) String() string {
	return "Owning organisation node changed: " + e.OwningOrgNodeID.String()
}
