package zenodo

import "time"

// Deposition represents a Zenodo deposition resource
type Deposition struct {
	ID        int                 `json:"id"`
	ConceptID string              `json:"conceptrecid,omitempty"`
	DOI       string              `json:"doi,omitempty"`
	DOIURL    string              `json:"doi_url,omitempty"`
	RecordID  int                 `json:"record_id,omitempty"`
	RecordURL string              `json:"record_url,omitempty"`
	Created   time.Time           `json:"created"`
	Modified  time.Time           `json:"modified"`
	Owner     int                 `json:"owner,omitempty"`
	State     DepositionState     `json:"state,omitempty"`
	Submitted bool                `json:"submitted,omitempty"`
	Title     string              `json:"title,omitempty"`
	Metadata  *DepositionMetadata `json:"metadata,omitempty"`
	Files     []DepositionFile    `json:"files,omitempty"`
	Links     *DepositionLinks    `json:"links,omitempty"`
}

// DepositionState represents the state of a deposition
type DepositionState string

const (
	StateInProgress  DepositionState = "inprogress"
	StateDone        DepositionState = "done"
	StateError       DepositionState = "error"
	StateUnsubmitted DepositionState = "unsubmitted"
)

// DepositionMetadata contains the metadata for a deposition
type DepositionMetadata struct {
	Title           string              `json:"title,omitempty"`
	UploadType      UploadType          `json:"upload_type,omitempty"`
	PublicationType PublicationType     `json:"publication_type,omitempty"`
	Description     string              `json:"description,omitempty"`
	Creators        []Creator           `json:"creators,omitempty"`
	AccessRight     AccessRight         `json:"access_right,omitempty"`
	License         string              `json:"license,omitempty"`
	DOI             string              `json:"doi,omitempty"`
	PrereserveDOI   *PrereserveDOI      `json:"prereserve_doi,omitempty"`
	Keywords        []string            `json:"keywords,omitempty"`
	Notes           string              `json:"notes,omitempty"`
	RelatedIDs      []RelatedIdentifier `json:"related_identifiers,omitempty"`
	Contributors    []Contributor       `json:"contributors,omitempty"`
	References      []string            `json:"references,omitempty"`
	Communities     []Community         `json:"communities,omitempty"`
	Grants          []Grant             `json:"grants,omitempty"`
	PublicationDate string              `json:"publication_date,omitempty"`
}

// PrereserveDOI contains pre-reserved DOI information
type PrereserveDOI struct {
	DOI   string `json:"doi"`
	RecID int    `json:"recid"`
}

// UploadType represents the type of upload
type UploadType string

const (
	UploadTypePublication    UploadType = "publication"
	UploadTypePoster         UploadType = "poster"
	UploadTypePresentation   UploadType = "presentation"
	UploadTypeDataset        UploadType = "dataset"
	UploadTypeImage          UploadType = "image"
	UploadTypeVideo          UploadType = "video"
	UploadTypeSoftware       UploadType = "software"
	UploadTypeLesson         UploadType = "lesson"
	UploadTypePhysicalObject UploadType = "physicalobject"
	UploadTypeOther          UploadType = "other"
)

// PublicationType represents the type of publication
type PublicationType string

const (
	PubTypeArticle               PublicationType = "article"
	PubTypeBook                  PublicationType = "book"
	PubTypeBookSection           PublicationType = "section"
	PubTypeConferencePaper       PublicationType = "conferencepaper"
	PubTypeDataManagementPlan    PublicationType = "datamanagementplan"
	PubTypePatent                PublicationType = "patent"
	PubTypePreprint              PublicationType = "preprint"
	PubTypeReport                PublicationType = "report"
	PubTypeSoftwareDocumentation PublicationType = "softwaredocumentation"
	PubTypeTechnicalNote         PublicationType = "technicalnote"
	PubTypeThesis                PublicationType = "thesis"
	PubTypeWorkingPaper          PublicationType = "workingpaper"
	PubTypeOther                 PublicationType = "other"
)

// AccessRight represents the access rights for a deposition
type AccessRight string

const (
	AccessOpen       AccessRight = "open"
	AccessEmbargoed  AccessRight = "embargoed"
	AccessRestricted AccessRight = "restricted"
	AccessClosed     AccessRight = "closed"
)

// Creator represents an author of the deposition
type Creator struct {
	Name        string `json:"name"` // Required. Format: "Family name, Given names"
	Affiliation string `json:"affiliation,omitempty"`
	ORCID       string `json:"orcid,omitempty"`
	GND         string `json:"gnd,omitempty"`
}

// Contributor represents a contributor to the deposition
type Contributor struct {
	Name        string          `json:"name"`
	Type        ContributorType `json:"type"`
	Affiliation string          `json:"affiliation,omitempty"`
	ORCID       string          `json:"orcid,omitempty"`
	GND         string          `json:"gnd,omitempty"`
}

// ContributorType represents the type of contributor
type ContributorType string

const (
	ContribContactPerson      ContributorType = "ContactPerson"
	ContribDataCollector      ContributorType = "DataCollector"
	ContribDataCurator        ContributorType = "DataCurator"
	ContribDataManager        ContributorType = "DataManager"
	ContribEditor             ContributorType = "Editor"
	ContribHostingInstitution ContributorType = "HostingInstitution"
	ContribProducer           ContributorType = "Producer"
	ContribProjectLeader      ContributorType = "ProjectLeader"
	ContribProjectManager     ContributorType = "ProjectManager"
	ContribProjectMember      ContributorType = "ProjectMember"
	ContribResearcher         ContributorType = "Researcher"
	ContribResearchGroup      ContributorType = "ResearchGroup"
	ContribRightsHolder       ContributorType = "RightsHolder"
	ContribSupervisor         ContributorType = "Supervisor"
	ContribSponsor            ContributorType = "Sponsor"
	ContribOther              ContributorType = "Other"
)

// RelatedIdentifier represents a related identifier
type RelatedIdentifier struct {
	Identifier   string `json:"identifier"`
	Relation     string `json:"relation"`
	ResourceType string `json:"resource_type,omitempty"`
}

// Community represents a Zenodo community
type Community struct {
	Identifier string `json:"identifier"`
}

// Grant represents a grant/funding source
type Grant struct {
	ID string `json:"id"`
}

// DepositionFile represents a file in a deposition
type DepositionFile struct {
	ID       string     `json:"id"`
	Filename string     `json:"filename"`
	Filesize int64      `json:"filesize"`
	Checksum string     `json:"checksum"`
	Links    *FileLinks `json:"links,omitempty"`
}

// FileLinks contains links for file operations
type FileLinks struct {
	Self     string `json:"self,omitempty"`
	Download string `json:"download,omitempty"`
}

// DepositionLinks contains HAL-style links for deposition operations
type DepositionLinks struct {
	Bucket          string `json:"bucket,omitempty"`
	Discard         string `json:"discard,omitempty"`
	Edit            string `json:"edit,omitempty"`
	Files           string `json:"files,omitempty"`
	HTML            string `json:"html,omitempty"`
	LatestDraft     string `json:"latest_draft,omitempty"`
	LatestDraftHTML string `json:"latest_draft_html,omitempty"`
	Publish         string `json:"publish,omitempty"`
	Self            string `json:"self,omitempty"`
	NewVersion      string `json:"newversion,omitempty"`
}

// TokenResponse represents the OAuth token response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
	UserID       string `json:"user_id,omitempty"`
}

// APIError represents an error response from the Zenodo API
type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return e.Message
}
