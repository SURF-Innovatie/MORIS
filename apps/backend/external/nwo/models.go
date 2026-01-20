package nwo

import "time"

// ProjectsResponse is the top-level response from the NWO Open API
type ProjectsResponse struct {
	Meta     Meta      `json:"meta"`
	Projects []Project `json:"projects"`
}

// Meta contains API metadata and pagination info
type Meta struct {
	APIType       string   `json:"api_type,omitempty"`
	Version       string   `json:"version,omitempty"`
	Licence       *Licence `json:"licence,omitempty"`
	Documentation string   `json:"documentation,omitempty"`
	ReleaseDate   string   `json:"release_date,omitempty"`
	Funder        string   `json:"funder,omitempty"`
	RORID         string   `json:"ror_id,omitempty"`
	Date          string   `json:"date,omitempty"`
	Count         int      `json:"count,omitempty"`
	PerPage       int      `json:"per_page,omitempty"`
	Pages         int      `json:"pages,omitempty"`
	Page          int      `json:"page,omitempty"`
}

// Licence contains licence information
type Licence struct {
	Name        string `json:"name,omitempty"`
	URL         string `json:"url,omitempty"`
	Description string `json:"decsription,omitempty"` // Note: typo exists in API spec
}

// Project represents a single NWO-subsidized project
type Project struct {
	ProjectID       string          `json:"project_id,omitempty"`
	GrantID         string          `json:"grant_id,omitempty"`
	ParentProjectID string          `json:"parent_project_id,omitempty"`
	Title           string          `json:"title,omitempty"`
	FundingSchemeID int64           `json:"funding_scheme_id,omitempty"`
	FundingScheme   string          `json:"funding_scheme,omitempty"`
	Department      string          `json:"department,omitempty"`
	SubDepartment   string          `json:"sub_department,omitempty"`
	StartDate       *time.Time      `json:"start_date,omitempty"`
	EndDate         *time.Time      `json:"end_date,omitempty"`
	ReportingYear   int             `json:"reporting_year,omitempty"`
	AwardAmount     int             `json:"award_amount,omitempty"`
	SummaryNL       string          `json:"summary_nl,omitempty"`
	SummaryEN       string          `json:"summary_en,omitempty"`
	DescriptionNL   string          `json:"Description_NL,omitempty"`
	DescriptionEN   string          `json:"Description_EN,omitempty"`
	SummaryUpdates  []SummaryUpdate `json:"summary_updates,omitempty"`
	ProjectMembers  []ProjectMember `json:"project_members,omitempty"`
	Products        []Product       `json:"Products,omitempty"`
}

// ProjectMember represents a member of a project team
type ProjectMember struct {
	Role              string `json:"role,omitempty"`
	MemberID          string `json:"member_id,omitempty"`
	ORCID             string `json:"orcid,omitempty"`
	LastName          string `json:"last_name,omitempty"`
	DegreePreNominal  string `json:"degree_pre_nominal,omitempty"`
	DegreePostNominal string `json:"degree_post_nominal,omitempty"`
	Initials          string `json:"initials,omitempty"`
	FirstName         string `json:"first_name,omitempty"`
	Prefix            string `json:"prefix,omitempty"`
	DAI               string `json:"dai,omitempty"`
	Organisation      string `json:"organisation,omitempty"`
	OrganisationID    string `json:"organisation_id,omitempty"`
	ROR               string `json:"ror,omitempty"`
	Active            string `json:"active,omitempty"`
}

// Product represents a publication or output from a project
type Product struct {
	ISBN          string   `json:"isbn,omitempty"`
	DOI           string   `json:"doi,omitempty"`
	Title         string   `json:"title,omitempty"`
	SubTitle      string   `json:"sub_title,omitempty"`
	Year          int      `json:"year,omitempty"`
	City          string   `json:"city,omitempty"`
	Edition       string   `json:"edition,omitempty"`
	Start         int      `json:"start,omitempty"`
	End           int      `json:"end,omitempty"`
	Type          string   `json:"type,omitempty"`
	URLOpenAccess string   `json:"url_open_access,omitempty"`
	Publisher     string   `json:"publisher,omitempty"`
	JournalTitle  string   `json:"journal_title,omitempty"`
	Authors       []Author `json:"authors,omitempty"`
}

// Author represents an author of a product
type Author struct {
	LastName          string `json:"last_name,omitempty"`
	DegreePreNominal  string `json:"degree_pre_nominal,omitempty"`
	DegreePostNominal string `json:"degree_post_nominal,omitempty"`
	Initials          string `json:"initials,omitempty"`
	FirstName         string `json:"first_name,omitempty"`
	Prefix            string `json:"prefix,omitempty"`
	DAI               string `json:"dai,omitempty"`
	Role              string `json:"role,omitempty"`
	IndexNumber       int    `json:"index_number,omitempty"`
}

// SummaryUpdate represents an update to the project summary
type SummaryUpdate struct {
	SubmissionDate *time.Time `json:"submission_date,omitempty"`
	UpdateEN       string     `json:"update_en,omitempty"`
	UpdateNL       string     `json:"update_nl,omitempty"`
}

// ProjectRole represents the possible roles for project members
type ProjectRole string

const (
	RoleMainApplicant ProjectRole = "Main Applicant"
	RoleCoApplicant   ProjectRole = "Co-applicant"
	RoleProjectLeader ProjectRole = "Project leader"
	RoleResearcher    ProjectRole = "Researcher"
)

// ExceptionResponse represents an error response from the API
type ExceptionResponse struct {
	Exception Exception `json:"exception"`
}

// Exception contains error details
type Exception struct {
	Timestamp string `json:"timestamp,omitempty"`
	Status    int    `json:"status,omitempty"`
	Error     string `json:"error,omitempty"`
	Message   string `json:"message,omitempty"`
}
