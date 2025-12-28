package events

import (
	"context"
	"errors"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

const OwningOrgNodeChangedType = "project.owning_org_node_changed"

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

type OwningOrgNodeChangedInput struct {
	OwningOrgNodeID uuid.UUID `json:"owning_org_node_id"`
}

func DecideOwningOrgNodeChanged(
	projectID uuid.UUID,
	actor uuid.UUID,
	cur *entities.Project,
	in OwningOrgNodeChangedInput,
	status Status,
) (Event, error) {
	if projectID == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	if in.OwningOrgNodeID == uuid.Nil {
		return nil, errors.New("organisation node id is required")
	}
	if cur == nil {
		return nil, errors.New("current project is required")
	}
	if cur.OwningOrgNodeID == in.OwningOrgNodeID {
		return nil, nil
	}

	return &OwningOrgNodeChanged{
		Base:            NewBase(projectID, actor, status),
		OwningOrgNodeID: in.OwningOrgNodeID,
	}, nil
}

func init() {
	RegisterMeta(EventMeta{
		Type:         OwningOrgNodeChangedType,
		FriendlyName: "Owning Organisation Node Change",
		CheckNotification: func(ctx context.Context, event Event, client *ent.Client) bool {
			return true
		},
	}, func() Event { return &OwningOrgNodeChanged{} })

	RegisterDecider[OwningOrgNodeChangedInput](OwningOrgNodeChangedType,
		func(ctx context.Context, projectID uuid.UUID, actor uuid.UUID, cur *entities.Project, in OwningOrgNodeChangedInput, status Status) (Event, error) {
			return DecideOwningOrgNodeChanged(projectID, actor, cur, in, status)
		})

	RegisterInputType(OwningOrgNodeChangedType, OwningOrgNodeChangedInput{})
}
