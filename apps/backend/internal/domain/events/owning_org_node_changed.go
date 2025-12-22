package events

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type OwningOrgNodeChanged struct {
	Base
	OwningOrgNodeID uuid.UUID `json:"owning_org_node_id"`
}

func (OwningOrgNodeChanged) isEvent()     {}
func (OwningOrgNodeChanged) Type() string { return OwningOrgNodeChangedType }
func (e OwningOrgNodeChanged) String() string {
	return "Owning organisation node changed"
}

func (e *OwningOrgNodeChanged) Apply(project *entities.Project) {
	project.OwningOrgNodeID = e.OwningOrgNodeID
}

func (e *OwningOrgNodeChanged) RelatedIDs() RelatedIDs {
	return RelatedIDs{OrgNodeID: &e.OwningOrgNodeID}
}

func init() {
	RegisterMeta(EventMeta{
		Type:         OwningOrgNodeChangedType,
		FriendlyName: "Owning Organisation Node Change",
		CheckNotification: func(ctx context.Context, event Event, client *ent.Client) bool {
			return true
		},
	}, func() Event { return &OwningOrgNodeChanged{} })
}
