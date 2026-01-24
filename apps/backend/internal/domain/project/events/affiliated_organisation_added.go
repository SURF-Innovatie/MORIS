package events

import (
	"context"
	"errors"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/internal/domain/project"
	"github.com/google/uuid"
)

const AffiliatedOrganisationAddedType = "project.affiliated_organisation_added"

type AffiliatedOrganisationAdded struct {
	Base
	AffiliatedOrganisationID uuid.UUID `json:"affiliated_organisation_id"`
}

func (AffiliatedOrganisationAdded) isEvent()     {}
func (AffiliatedOrganisationAdded) Type() string { return AffiliatedOrganisationAddedType }
func (e AffiliatedOrganisationAdded) String() string {
	return fmt.Sprintf("Affiliated Organisation added: %s", e.AffiliatedOrganisationID)
}

func (e *AffiliatedOrganisationAdded) Apply(project *project.Project) {
	project.AffiliatedOrganisationIDs = append(project.AffiliatedOrganisationIDs, e.AffiliatedOrganisationID)
}

func (e *AffiliatedOrganisationAdded) NotificationMessage() string {
	return "A new affiliated organisation has been added to the project."
}

type AffiliatedOrganisationAddedInput struct {
	AffiliatedOrganisationID uuid.UUID `json:"affiliated_organisation_id"`
}

func DecideAffiliatedOrganisationAdded(
	projectID uuid.UUID,
	actor uuid.UUID,
	cur *project.Project,
	in AffiliatedOrganisationAddedInput,
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

	for _, x := range cur.AffiliatedOrganisationIDs {
		if x == in.AffiliatedOrganisationID {
			return nil, fmt.Errorf("affiliated organisation %s already exists in project %s", in.AffiliatedOrganisationID, cur.Id)
		}
	}

	base := NewBase(projectID, actor, status)
	base.FriendlyNameStr = AffiliatedOrganisationAddedMeta.FriendlyName

	return &AffiliatedOrganisationAdded{
		Base:                     base,
		AffiliatedOrganisationID: in.AffiliatedOrganisationID,
	}, nil
}

var AffiliatedOrganisationAddedMeta = EventMeta{
	Type:         AffiliatedOrganisationAddedType,
	FriendlyName: "Affiliated Organisation Addition",
}

func init() {
	RegisterMeta(AffiliatedOrganisationAddedMeta, func() Event {
		return &AffiliatedOrganisationAdded{
			Base: Base{FriendlyNameStr: AffiliatedOrganisationAddedMeta.FriendlyName},
		}
	})

	RegisterDecider(AffiliatedOrganisationAddedType,
		func(ctx context.Context, projectID uuid.UUID, actor uuid.UUID, cur *project.Project, in AffiliatedOrganisationAddedInput, status Status) (Event, error) {
			return DecideAffiliatedOrganisationAdded(projectID, actor, cur, in, status)
		})

	RegisterInputType(AffiliatedOrganisationAddedType, AffiliatedOrganisationAddedInput{})
}
