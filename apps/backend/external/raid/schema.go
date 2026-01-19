package raid

// Schema URIs for RAiD metadata
// See: https://metadata.raid.org/en/latest/

// Title Type Schema
// See: https://vocabulary.raid.org/title.type.schema/376
const (
	TitleTypePrimaryURI     = "https://vocabulary.raid.org/title.type.schema/5"
	TitleTypeAlternativeURI = "https://vocabulary.raid.org/title.type.schema/4"
	TitleTypeSchemaURI      = "https://vocabulary.raid.org/title.type.schema/376"
)

// Description Type Schema
// See: https://vocabulary.raid.org/description.type.schema/329
const (
	DescriptionTypePrimaryURI      = "https://vocabulary.raid.org/description.type.schema/318"
	DescriptionTypeAlternativeURI  = "https://vocabulary.raid.org/description.type.schema/319"
	DescriptionTypeSignificanceURI = "https://vocabulary.raid.org/description.type.schema/320"
	DescriptionTypeMethodsURI      = "https://vocabulary.raid.org/description.type.schema/321"
	DescriptionTypeSchemaURI       = "https://vocabulary.raid.org/description.type.schema/329"
)

// Access Type Schema
// See: https://vocabulary.raid.org/access.type.schema/289
const (
	AccessTypeOpenURI      = "https://vocabulary.raid.org/access.type.schema/238"
	AccessTypeEmbargoedURI = "https://vocabulary.raid.org/access.type.schema/239"
	AccessTypeSchemaURI    = "https://vocabulary.raid.org/access.type.schema/289"
	// COAR access rights (alternative schema)
	AccessTypeCOAROpenURI   = "https://vocabularies.coar-repositories.org/access_rights/c_abf2/"
	AccessTypeCOARSchemaURI = "https://vocabularies.coar-repositories.org/access_rights/"
)

// Language Schema (ISO 639-3)
const (
	LanguageSchemaURI = "https://www.iso.org/standard/39534.html"
	LanguageEnglish   = "eng"
	LanguageDutch     = "nld"
)

// Contributor Schema (ORCID)
const (
	ContributorSchemaORCID = "https://orcid.org/"
)

// Contributor Position Schema
// See: https://vocabulary.raid.org/contributor.position.schema/305
const (
	ContributorPositionLeaderURI       = "https://vocabulary.raid.org/contributor.position.schema/306"
	ContributorPositionOtherURI        = "https://vocabulary.raid.org/contributor.position.schema/307"
	ContributorPositionCoInvestigatURI = "https://vocabulary.raid.org/contributor.position.schema/308"
	ContributorPositionContactURI      = "https://vocabulary.raid.org/contributor.position.schema/309"
	ContributorPositionSchemaURI       = "https://vocabulary.raid.org/contributor.position.schema/305"
)

// Contributor Role Schema (CRediT - https://credit.niso.org/)
const (
	ContributorRoleSchemaURI            = "https://credit.niso.org/"
	ContributorRoleConceptualizationURI = "https://credit.niso.org/contributor-roles/conceptualization/"
	ContributorRoleDataCurationURI      = "https://credit.niso.org/contributor-roles/data-curation/"
	ContributorRoleFormalAnalysisURI    = "https://credit.niso.org/contributor-roles/formal-analysis/"
	ContributorRoleFundingAcquisURI     = "https://credit.niso.org/contributor-roles/funding-acquisition/"
	ContributorRoleInvestigationURI     = "https://credit.niso.org/contributor-roles/investigation/"
	ContributorRoleMethodologyURI       = "https://credit.niso.org/contributor-roles/methodology/"
	ContributorRoleProjectAdminURI      = "https://credit.niso.org/contributor-roles/project-administration/"
	ContributorRoleResourcesURI         = "https://credit.niso.org/contributor-roles/resources/"
	ContributorRoleSoftwareURI          = "https://credit.niso.org/contributor-roles/software/"
	ContributorRoleSupervisionURI       = "https://credit.niso.org/contributor-roles/supervision/"
	ContributorRoleValidationURI        = "https://credit.niso.org/contributor-roles/validation/"
	ContributorRoleVisualizationURI     = "https://credit.niso.org/contributor-roles/visualization/"
	ContributorRoleWritingOrigURI       = "https://credit.niso.org/contributor-roles/writing-original-draft/"
	ContributorRoleWritingReviewURI     = "https://credit.niso.org/contributor-roles/writing-review-editing/"
)

// Organisation Schema (ROR)
const (
	OrganisationSchemaROR = "https://ror.org/"
)

// Organisation Role Schema
// See: https://vocabulary.raid.org/organisation.role.schema/359
const (
	OrganisationRoleLeadResearchURI = "https://vocabulary.raid.org/organisation.role.schema/182"
	OrganisationRolePartnerURI      = "https://vocabulary.raid.org/organisation.role.schema/183"
	OrganisationRoleFunderURI       = "https://vocabulary.raid.org/organisation.role.schema/184"
	OrganisationRoleFacilityURI     = "https://vocabulary.raid.org/organisation.role.schema/185"
	OrganisationRoleOtherURI        = "https://vocabulary.raid.org/organisation.role.schema/186"
	OrganisationRoleContractorURI   = "https://vocabulary.raid.org/organisation.role.schema/187"
	OrganisationRoleSchemaURI       = "https://vocabulary.raid.org/organisation.role.schema/359"
)

// Related Object Type Schema
// See: https://vocabulary.raid.org/relatedObject.type.schema/329
const (
	RelatedObjectTypeSchemaURI = "https://vocabulary.raid.org/relatedObject.type.schema/329"
)

// Related Object Category Schema
const (
	RelatedObjectCategorySchemaURI = "https://vocabulary.raid.org/relatedObject.category.schema/"
)

// RAiD Identifier Schema
const (
	RAiDIdentifierSchemaURI = "https://raid.org/"
)

// Constraints
const (
	MaxTitleLength       = 100
	MaxDescriptionLength = 1000
)
