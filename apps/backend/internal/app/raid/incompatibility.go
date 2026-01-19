package raid

import "github.com/google/uuid"

// IncompatibilityType represents a type of RAiD compatibility issue.
type IncompatibilityType string

// Incompatibility types matching the C# RAiDIncompatibilityType enum
const (
	// Title incompatibilities
	NoActivePrimaryTitle       IncompatibilityType = "no_active_primary_title"
	MultipleActivePrimaryTitle IncompatibilityType = "multiple_active_primary_title"
	ProjectTitleTooLong        IncompatibilityType = "project_title_too_long"

	// Description incompatibilities
	NoPrimaryDescription       IncompatibilityType = "no_primary_description"
	MultiplePrimaryDescription IncompatibilityType = "multiple_primary_descriptions"
	ProjectDescriptionTooLong  IncompatibilityType = "project_description_too_long"

	// Contributor incompatibilities
	NoContributors                  IncompatibilityType = "no_contributors"
	ContributorWithoutOrcid         IncompatibilityType = "contributor_without_orcid"
	OverlappingContributorPositions IncompatibilityType = "overlapping_contributor_positions"
	ContributorWithoutPosition      IncompatibilityType = "contributor_without_position"
	NoProjectLeader                 IncompatibilityType = "no_project_leader"
	NoProjectContact                IncompatibilityType = "no_project_contact"

	// Organisation incompatibilities
	OrganisationWithoutRor           IncompatibilityType = "organisation_without_ror"
	OverlappingOrganisationRoles     IncompatibilityType = "overlapping_organisation_roles"
	NoLeadResearchOrganisation       IncompatibilityType = "no_lead_research_organisation"
	MultipleLeadResearchOrganisation IncompatibilityType = "multiple_lead_research_organisation"

	// Product incompatibilities
	NoProductCategory IncompatibilityType = "no_product_category"

	// Language incompatibilities
	InvalidTitleLanguage       IncompatibilityType = "invalid_title_language"
	InvalidDescriptionLanguage IncompatibilityType = "invalid_description_language"
)

// Incompatibility represents a single RAiD compatibility issue.
type Incompatibility struct {
	Type     IncompatibilityType `json:"type"`
	ObjectID uuid.UUID           `json:"objectId,omitempty"`
	Message  string              `json:"message,omitempty"`
}

// NewIncompatibility creates a new Incompatibility with the given type.
func NewIncompatibility(t IncompatibilityType) Incompatibility {
	return Incompatibility{Type: t}
}

// WithObjectID adds an object ID to the incompatibility.
func (i Incompatibility) WithObjectID(id uuid.UUID) Incompatibility {
	i.ObjectID = id
	return i
}

// WithMessage adds a descriptive message to the incompatibility.
func (i Incompatibility) WithMessage(msg string) Incompatibility {
	i.Message = msg
	return i
}
