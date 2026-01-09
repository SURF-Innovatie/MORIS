package raid

// RAiDId represents the identifier of a RAiD.
type RAiDId struct {
	IdValue            string                 `json:"id"`
	SchemaUri          string                 `json:"schemaUri"`
	RegistrationAgency RAiDRegistrationAgency `json:"registrationAgency"`
	Owner              RAiDOwner              `json:"owner"`
	RaidAgencyUrl      *string                `json:"raidAgencyUrl,omitempty"`
	License            string                 `json:"license"`
	Version            int                    `json:"version"`
}

// RAiDRegistrationAgency represents the agency that registered the RAiD.
type RAiDRegistrationAgency struct {
	Id        string `json:"id"`
	SchemaUri string `json:"schemaUri"`
}

// RAiDOwner represents the owner of the RAiD.
type RAiDOwner struct {
	Id        string `json:"id"`
	SchemaUri string `json:"schemaUri"`
	ServicePoint *string `json:"servicePoint,omitempty"`
}

// RAiDTitle represents a title of the RAiD.
type RAiDTitle struct {
	Text      string        `json:"text"`
	Type      RAiDTitleType `json:"type"`
	StartDate string        `json:"startDate"` // YYYY-MM-DD
	EndDate   *string       `json:"endDate,omitempty"` // YYYY-MM-DD
	Language  *RAiDLanguage `json:"language,omitempty"`
}

// RAiDTitleType represents the type of title.
type RAiDTitleType struct {
	Id        string `json:"id"`
	SchemaUri string `json:"schemaUri"`
}

// RAiDLanguage represents a language.
type RAiDLanguage struct {
	Id        string `json:"id"`
	SchemaUri string `json:"schemaUri"`
}

// RAiDDate represents the date coverage of the RAiD.
type RAiDDate struct {
	StartDate string  `json:"startDate"` // YYYY-MM-DD
	EndDate   *string `json:"endDate,omitempty"` // YYYY-MM-DD
}

// RAiDDescription represents a description of the RAiD.
type RAiDDescription struct {
	Text     string              `json:"text"`
	Type     RAiDDescriptionType `json:"type"`
	Language *RAiDLanguage       `json:"language,omitempty"`
}

// RAiDDescriptionType represents the type of description.
type RAiDDescriptionType struct {
	Id        string `json:"id"`
	SchemaUri string `json:"schemaUri"`
}

// RAiDAccess represents the access conditions of the RAiD.
type RAiDAccess struct {
	Type          RAiDAccessType       `json:"type"`
	Statement     *RAiDAccessStatement `json:"statement,omitempty"`
	EmbargoExpiry *string              `json:"embargoExpiry,omitempty"` // YYYY-MM-DD
}

// RAiDAccessType represents the type of access.
type RAiDAccessType struct {
	Id        string `json:"id"`
	SchemaUri string `json:"schemaUri"`
}

// RAiDAccessStatement represents the access statement.
type RAiDAccessStatement struct {
	Text     string        `json:"text"`
	Language *RAiDLanguage `json:"language,omitempty"`
}

// RAiDAlternateUrl represents an alternate URL for the RAiD.
type RAiDAlternateUrl struct {
	Url string `json:"url"`
}

// RAiDContributor represents a contributor to the RAiD.
type RAiDContributor struct {
	Id            *string                 `json:"id,omitempty"` // Note: Some bad data in RAiD db might make this missing
	SchemaUri     string                  `json:"schemaUri"`
	Status        *string                 `json:"status,omitempty"`
	StatusMessage *string                 `json:"statusMessage,omitempty"`
	Email         *string                 `json:"email,omitempty"`
	Uuid          *string                 `json:"uuid,omitempty"`
	Position      []RAiDContributorPosition `json:"position"`
	Role          []RAiDContributorRole   `json:"role"`
	Leader        *bool                   `json:"leader,omitempty"`
	Contact       *bool                   `json:"contact,omitempty"`
}

// RAiDContributorPosition represents the position of a contributor.
type RAiDContributorPosition struct {
	Id        string      `json:"id"`
	SchemaUri string      `json:"schemaUri"`
	StartDate string      `json:"startDate"` // YYYY-MM-DD
	EndDate   *string     `json:"endDate,omitempty"` // YYYY-MM-DD
}

// RAiDContributorRole represents the role of a contributor.
type RAiDContributorRole struct {
	Id        string `json:"id"`
	SchemaUri string `json:"schemaUri"`
}

// RAiDOrganisation represents an organisation linked to the RAiD.
type RAiDOrganisation struct {
	Id        string                 `json:"id"`
	SchemaUri string                 `json:"schemaUri"`
	Role      []RAiDOrganisationRole `json:"role"`
}

// RAiDOrganisationRole represents the role of an organisation.
type RAiDOrganisationRole struct {
	Id        string      `json:"id"`
	SchemaUri string      `json:"schemaUri"`
	StartDate string      `json:"startDate"` // YYYY-MM-DD
	EndDate   *string     `json:"endDate,omitempty"` // YYYY-MM-DD
}

// RAiDSubject represents a subject of the RAiD.
type RAiDSubject struct {
	Id        string             `json:"id"`
	SchemaUri string             `json:"schemaUri"`
	Keyword   []RAiDSubjectKeyword `json:"keyword,omitempty"`
}

// RAiDSubjectKeyword represents a keyword for a subject.
type RAiDSubjectKeyword struct {
	Text     string        `json:"text"`
	Language *RAiDLanguage `json:"language,omitempty"`
}

// RAiDRelatedRaid represents a related RAiD.
type RAiDRelatedRaid struct {
	Id        string               `json:"id"`
	Type      RAiDRelatedRaidType  `json:"type"`
	Title     *string              `json:"title,omitempty"`
}

// RAiDRelatedRaidType represents the type of related RAiD relationship.
type RAiDRelatedRaidType struct {
	Id        string `json:"id"`
	SchemaUri string `json:"schemaUri"`
}

// RAiDRelatedObject represents a related object.
type RAiDRelatedObject struct {
	Id        string                    `json:"id"`
	SchemaUri *string                   `json:"schemaUri,omitempty"`
	Type      RAiDRelatedObjectType     `json:"type"`
	Category  []RAiDRelatedObjectCategory `json:"category,omitempty"`
	Title     *string                   `json:"title,omitempty"`
}

// RAiDRelatedObjectType represents the type of related object.
type RAiDRelatedObjectType struct {
	Id        string `json:"id"`
	SchemaUri string `json:"schemaUri"`
}

// RAiDRelatedObjectCategory represents the category of related object.
type RAiDRelatedObjectCategory struct {
	Id        string `json:"id"`
	SchemaUri string `json:"schemaUri"`
}

// RAiDAlternateIdentifier represents an alternate identifier.
type RAiDAlternateIdentifier struct {
	Id     string `json:"id"`
	Type   string `json:"type"`
}

// RAiDSpatialCoverage represents spatial coverage data.
type RAiDSpatialCoverage struct {
	Id        string                  `json:"id"`
	SchemaUri string                  `json:"schemaUri"`
	Place     []RAiDSpatialCoveragePlace `json:"place,omitempty"`
	Language  *RAiDLanguage           `json:"language,omitempty"`
}

// RAiDSpatialCoveragePlace represents a place in spatial coverage.
type RAiDSpatialCoveragePlace struct {
	Text     string        `json:"text"`
	Language *RAiDLanguage `json:"language,omitempty"`
}

// RAiDTraditionalKnowledgeLabel represents traditional knowledge labels.
type RAiDTraditionalKnowledgeLabel struct {
	Id        string `json:"id,omitempty"`
	SchemaUri string `json:"schemaUri,omitempty"`
}

// RAiDMetadata represents metadata about the RAiD record itself.
type RAiDMetadata struct {
	Access             string `json:"access"`
	IdProvider         string `json:"idProvider"`
	IdProviderClient   string `json:"idProviderClient"`
}

// RAiDDto represents the full RAiD object.
type RAiDDto struct {
	TraditionalKnowledgeLabel []RAiDTraditionalKnowledgeLabel `json:"traditionalKnowledgeLabel,omitempty"`
	Metadata                  *RAiDMetadata                   `json:"metadata,omitempty"`
	Identifier                RAiDId                          `json:"identifier"`
	Title                     []RAiDTitle                     `json:"title,omitempty"`
	Date                      *RAiDDate                       `json:"date,omitempty"`
	Description               []RAiDDescription               `json:"description,omitempty"`
	Access                    RAiDAccess                      `json:"access"`
	AlternateUrl              []RAiDAlternateUrl              `json:"alternateUrl,omitempty"`
	Contributor               []RAiDContributor               `json:"contributor,omitempty"`
	Organisation              []RAiDOrganisation              `json:"organisation,omitempty"`
	Subject                   []RAiDSubject                   `json:"subject,omitempty"`
	RelatedRaid               []RAiDRelatedRaid               `json:"relatedRaid,omitempty"`
	RelatedObject             []RAiDRelatedObject             `json:"relatedObject,omitempty"`
	AlternateIdentifier       []RAiDAlternateIdentifier       `json:"alternateIdentifier,omitempty"`
	SpatialCoverage           []RAiDSpatialCoverage           `json:"spatialCoverage,omitempty"`
}

// RAiDCreateRequest represents the payload for creating a new RAiD.
type RAiDCreateRequest struct {
	Metadata            *RAiDMetadata                   `json:"metadata,omitempty"`
	Identifier          *RAiDId                         `json:"identifier,omitempty"` // For minting, usually we don't send ID, or we send partial? C# says nullable.
	Title               []RAiDTitle                     `json:"title,omitempty"`
	Date                *RAiDDate                       `json:"date,omitempty"`
	Description         []RAiDDescription               `json:"description,omitempty"`
	Access              RAiDAccess                      `json:"access"`
	AlternateUrl        []RAiDAlternateUrl              `json:"alternateUrl,omitempty"`
	Contributor         []RAiDContributor               `json:"contributor,omitempty"`
	Organisation        []RAiDOrganisation              `json:"organisation,omitempty"`
	Subject             []RAiDSubject                   `json:"subject,omitempty"`
	RelatedRaid         []RAiDRelatedRaid               `json:"relatedRaid,omitempty"`
	RelatedObject       []RAiDRelatedObject             `json:"relatedObject,omitempty"`
	AlternateIdentifier []RAiDAlternateIdentifier       `json:"alternateIdentifier,omitempty"`
	SpatialCoverage     []RAiDSpatialCoverage           `json:"spatialCoverage,omitempty"`
}

// RAiDUpdateRequest represents the payload for updating a RAiD.
type RAiDUpdateRequest struct {
	Title               []RAiDTitle                     `json:"title,omitempty"`
	Date                *RAiDDate                       `json:"date,omitempty"`
	Description         []RAiDDescription               `json:"description,omitempty"`
	Access              RAiDAccess                      `json:"access"`
	AlternateUrl        []RAiDAlternateUrl              `json:"alternateUrl,omitempty"`
	Contributor         []RAiDContributor               `json:"contributor,omitempty"`
	Organisation        []RAiDOrganisation              `json:"organisation,omitempty"`
	Subject             []RAiDSubject                   `json:"subject,omitempty"`
	RelatedRaid         []RAiDRelatedRaid               `json:"relatedRaid,omitempty"`
	RelatedObject       []RAiDRelatedObject             `json:"relatedObject,omitempty"`
	AlternateIdentifier []RAiDAlternateIdentifier       `json:"alternateIdentifier,omitempty"`
	SpatialCoverage     []RAiDSpatialCoverage           `json:"spatialCoverage,omitempty"`
	Identifier          RAiDId                          `json:"identifier"`
}

// RAiDAuthResponse represents the response from the authentication endpoint.
type RAiDAuthResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	IdToken          string `json:"id_token"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}
