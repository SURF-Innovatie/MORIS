package events

import (
	"context"
	"errors"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/internal/domain/project"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

const AffiliatedOrganisationRemovedType = "project.affiliated_organisation_removed"

type AffiliatedOrganisationRemoved struct {
	Base
	AffiliatedOrganisationID uuid.UUID `json:"affiliated_organisation_id"`
}

func (AffiliatedOrganisationRemoved) isEvent()     {}
func (AffiliatedOrganisationRemoved) Type() string { return AffiliatedOrganisationRemovedType }
func (e AffiliatedOrganisationRemoved) String() string {
	return fmt.Sprintf("Affiliated Organisation removed: %s", e.AffiliatedOrganisationID)
}

func (e *AffiliatedOrganisationRemoved) Apply(project *project.Project) {
	project.AffiliatedOrganisationIDs = lo.Filter(project.AffiliatedOrganisationIDs, func(id uuid.UUID, _ int) bool {
		return id != e.AffiliatedOrganisationID
	})
}

func (e *AffiliatedOrganisationRemoved) NotificationMessage() string {
	return "An affiliated organisation has been removed from the project."
}

type AffiliatedOrganisationRemovedInput struct {
	AffiliatedOrganisationID uuid.UUID `json:"affiliated_organisation_id"`
}

func DecideAffiliatedOrganisationRemoved(
	projectID uuid.UUID,
	actor uuid.UUID,
	cur *project.Project,
	in AffiliatedOrganisationRemovedInput,
	status Status,
) (Event, error) {
	if projectID == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	if in.AffiliatedOrganisationID == uuid.Nil {
		return nil, errors.New("affiliated organisation id is required")
	}
	if cur == nil {
		return nil, errors.New("current project is required")
	}

	exists := false
	for _, x := range cur.AffiliatedOrganisationIDs {
		if x == in.AffiliatedOrganisationID {
			exists = true
			break
		}
	}
	if !exists {
		return nil, fmt.Errorf("affiliated organisation %s does not exist in project %s", in.AffiliatedOrganisationID, cur.Id)
	}

	base := NewBase(projectID, actor, status)
	base.FriendlyNameStr = AffiliatedOrganisationRemovedMeta.FriendlyName

	return &AffiliatedOrganisationRemoved{
		Base:                     base,
		AffiliatedOrganisationID: in.AffiliatedOrganisationID,
	}, nil
}

var AffiliatedOrganisationRemovedMeta = EventMeta{
	Type:         AffiliatedOrganisationRemovedType,
	FriendlyName: "Affiliated Organisation Removal",
}

func init() {
	RegisterMeta(AffiliatedOrganisationRemovedMeta, func() Event {
		return &AffiliatedOrganisationRemoved{
			Base: Base{FriendlyNameStr: AffiliatedOrganisationRemovedMeta.FriendlyName},
		}
	})

	RegisterDecider(AffiliatedOrganisationRemovedType,
		func(ctx context.Context, projectID uuid.UUID, actor uuid.UUID, cur *project.Project, in AffiliatedOrganisationRemovedInput, status Status) (Event, error) {
			return DecideAffiliatedOrganisationRemoved(projectID, actor, cur, in, status)
		})

	RegisterInputType(AffiliatedOrganisationRemovedType, AffiliatedOrganisationRemovedInput{})
}
