package crossref

import "time"

// Work represents a scholarly work from Crossref API
type Work struct {
	Indexed             *Indexed         `json:"indexed,omitempty"`
	ReferenceCount      int              `json:"reference-count,omitempty"`
	Publisher           string           `json:"publisher,omitempty"`
	Issue               string           `json:"issue,omitempty"`
	License             []License        `json:"license,omitempty"`
	ContentDomain       *ContentDomain   `json:"content-domain,omitempty"`
	ShortContainerTitle []string         `json:"short-container-title,omitempty"`
	PublishedPrint      *PublishedPrint  `json:"published-print,omitempty"`
	DOI                 string           `json:"DOI,omitempty"`
	Type                string           `json:"type,omitempty"`
	Created             *Created         `json:"created,omitempty"`
	Page                string           `json:"page,omitempty"`
	Source              string           `json:"source,omitempty"`
	IsReferencedByCount int              `json:"is-referenced-by-count,omitempty"`
	Title               []string         `json:"title,omitempty"`
	Prefix              string           `json:"prefix,omitempty"`
	Volume              string           `json:"volume,omitempty"`
	Author              []Author         `json:"author,omitempty"`
	Member              string           `json:"member,omitempty"`
	Reference           []Reference      `json:"reference,omitempty"`
	ContainerTitle      []string         `json:"container-title,omitempty"`
	Language            string           `json:"language,omitempty"`
	Link                []Link           `json:"link,omitempty"`
	Deposited           *Deposited       `json:"deposited,omitempty"`
	Score               float64          `json:"score,omitempty"`
	Resource            *WorkResource    `json:"resource,omitempty"`
	Issued              *Issued          `json:"issued,omitempty"`
	ReferencesCount     int              `json:"references-count,omitempty"`
	JournalIssue        *JournalIssue    `json:"journal-issue,omitempty"`
	AlternativeID       []string         `json:"alternative-id,omitempty"`
	URL                 string           `json:"URL,omitempty"`
	ISSN                []string         `json:"ISSN,omitempty"`
	ISSNType            []ISSNType       `json:"issn-type,omitempty"`
	Published           *Published       `json:"published,omitempty"`
	PublishedOnline     *PublishedOnline `json:"published-online,omitempty"`
	Abstract            string           `json:"abstract,omitempty"`
	Subtitle            []string         `json:"subtitle,omitempty"`
	Archive             []string         `json:"archive,omitempty"`
	UpdatePolicy        string           `json:"update-policy,omitempty"`
	Assertion           []Assertion      `json:"assertion,omitempty"`
}

// Author represents an author of a work
type Author struct {
	Given       string              `json:"given,omitempty"`
	Family      string              `json:"family"`
	Sequence    string              `json:"sequence"`
	Affiliation []AuthorAffiliation `json:"affiliation"`
}

// AuthorAffiliation represents an author's affiliation
type AuthorAffiliation struct {
	Name string `json:"name,omitempty"`
}

// Journal represents a journal from Crossref API
type Journal struct {
	LastStatusCheckTime uint64                   `json:"last-status-check-time,omitempty"`
	Counts              *JournalCounts           `json:"counts"`
	Breakdowns          *JournalBreakdowns       `json:"breakdowns"`
	Publisher           string                   `json:"publisher"`
	Coverage            *CombinedJournalCoverage `json:"coverage"`
	Title               string                   `json:"title"`
	Subjects            []any                    `json:"subjects"`
	CoverageType        *JournalCoverageType     `json:"coverage-type"`
	Flags               *JournalFlags            `json:"flags"`
	ISSN                []string                 `json:"ISSN"`
	ISSNType            []ISSNType               `json:"issn-type"`
}

// JournalCounts represents journal metadata counts
type JournalCounts struct {
	Total        int `json:"total-dois,omitempty"`
	CurrentDOIs  int `json:"current-dois,omitempty"`
	BackfileDOIs int `json:"backfile-dois,omitempty"`
}

// JournalBreakdowns represents journal breakdowns
type JournalBreakdowns struct {
	DOIsByIssuedYear [][]int `json:"dois-by-issued-year,omitempty"`
}

// CombinedJournalCoverage represents journal coverage data
type CombinedJournalCoverage struct {
	AffiliationsCurrent        float64 `json:"affiliations-current,omitempty"`
	SimilarityCheckingCurrent  float64 `json:"similarity-checking-current,omitempty"`
	DescriptionsCurrent        float64 `json:"descriptions-current,omitempty"`
	RorIdsCurrent              float64 `json:"ror-ids-current,omitempty"`
	FundersBackfile            float64 `json:"funders-backfile,omitempty"`
	LicensesBackfile           float64 `json:"licenses-backfile,omitempty"`
	FundersCurrent             float64 `json:"funders-current,omitempty"`
	AffiliationsBackfile       float64 `json:"affiliations-backfile,omitempty"`
	ResourceLinksBackfile      float64 `json:"resource-links-backfile,omitempty"`
	ORCIDsBackfile             float64 `json:"orcids-backfile,omitempty"`
	UpdatePoliciesCurrent      float64 `json:"update-policies-current,omitempty"`
	RorIdsBackfile             float64 `json:"ror-ids-backfile,omitempty"`
	ORCIDsCurrent              float64 `json:"orcids-current,omitempty"`
	SimilarityCheckingBackfile float64 `json:"similarity-checking-backfile,omitempty"`
	ReferencesBackfile         float64 `json:"references-backfile,omitempty"`
	DescriptionsBackfile       float64 `json:"descriptions-backfile,omitempty"`
	AwardNumbersBackfile       float64 `json:"award-numbers-backfile,omitempty"`
	UpdatePoliciesBackfile     float64 `json:"update-policies-backfile,omitempty"`
	LicensesCurrent            float64 `json:"licenses-current,omitempty"`
	AwardNumbersCurrent        float64 `json:"award-numbers-current,omitempty"`
	AbstractsBackfile          float64 `json:"abstracts-backfile,omitempty"`
	ResourceLinksCurrent       float64 `json:"resource-links-current,omitempty"`
	AbstractsCurrent           float64 `json:"abstracts-current,omitempty"`
	ReferencesCurrent          float64 `json:"references-current,omitempty"`
}

// JournalCoverageType represents journal coverage type
type JournalCoverageType struct {
	All      *JournalCoverage `json:"all,omitempty"`
	Current  *JournalCoverage `json:"current,omitempty"`
	Backfile *JournalCoverage `json:"backfile,omitempty"`
}

// JournalCoverage represents detailed journal coverage
type JournalCoverage struct {
	LastStatusCheckTime uint64  `json:"last-status-check-time,omitempty"`
	Affiliations        float64 `json:"affiliations,omitempty"`
	Abstracts           float64 `json:"abstracts,omitempty"`
	ORCIDs              float64 `json:"orcids,omitempty"`
	Licenses            float64 `json:"licenses,omitempty"`
	References          float64 `json:"references,omitempty"`
	Funders             float64 `json:"funders,omitempty"`
	SimilarityChecking  float64 `json:"similarity-checking,omitempty"`
	AwardNumbers        float64 `json:"award-numbers,omitempty"`
	RorIds              float64 `json:"ror-ids,omitempty"`
	UpdatePolicies      float64 `json:"update-policies,omitempty"`
	ResourceLinks       float64 `json:"resource-links,omitempty"`
	Descriptions        float64 `json:"descriptions,omitempty"`
}

// JournalFlags represents journal flags
type JournalFlags struct {
	DepositsAbstractsCurrent           bool `json:"deposits-abstracts-current,omitempty"`
	DepositsORCIDsCurrent              bool `json:"deposits-orcids-current,omitempty"`
	Deposits                           bool `json:"deposits,omitempty"`
	DepositsAffiliationsBackfile       bool `json:"deposits-affiliations-backfile,omitempty"`
	DepositsUpdatePoliciesBackfile     bool `json:"deposits-update-policies-backfile,omitempty"`
	DepositsSimilarityCheckingBackfile bool `json:"deposits-similarity-checking-backfile,omitempty"`
	DepositsAwardNumbersCurrent        bool `json:"deposits-award-numbers-current,omitempty"`
	DepositsResourceLinksCurrent       bool `json:"deposits-resource-links-current,omitempty"`
	DepositsRorIdsCurrent              bool `json:"deposits-ror-ids-current,omitempty"`
	DepositsArticles                   bool `json:"deposits-articles,omitempty"`
	DepositsAffiliationsCurrent        bool `json:"deposits-affiliations-current,omitempty"`
	DepositsFundersCurrent             bool `json:"deposits-funders-current,omitempty"`
	DepositsReferencesBackfile         bool `json:"deposits-references-backfile,omitempty"`
	DepositsRorIdsBackfile             bool `json:"deposits-ror-ids-backfile,omitempty"`
	DepositsAbstractsBackfile          bool `json:"deposits-abstracts-backfile,omitempty"`
	DepositsLicensesBackfile           bool `json:"deposits-licenses-backfile,omitempty"`
	DepositsAwardNumbersBackfile       bool `json:"deposits-award-numbers-backfile,omitempty"`
	DepositsDescriptionsCurrent        bool `json:"deposits-descriptions-current,omitempty"`
	DepositsReferencesCurrent          bool `json:"deposits-references-current,omitempty"`
	DepositsResourceLinksBackfile      bool `json:"deposits-resource-links-backfile,omitempty"`
	DepositsDescriptionsBackfile       bool `json:"deposits-descriptions-backfile,omitempty"`
	DepositsORCIDsBackfile             bool `json:"deposits-orcids-backfile,omitempty"`
	DepositsFundersBackfile            bool `json:"deposits-funders-backfile,omitempty"`
	DepositsUpdatePoliciesCurrent      bool `json:"deposits-update-policies-current,omitempty"`
	DepositsSimilarityCheckingCurrent  bool `json:"deposits-similarity-checking-current,omitempty"`
	DepositsLicensesCurrent            bool `json:"deposits-licenses-current,omitempty"`
}

// License represents a license
type License struct {
	URL            string     `json:"URL,omitempty"`
	Start          *DateParts `json:"start,omitempty"`
	DelayInDays    int        `json:"delay-in-days,omitempty"`
	ContentVersion string     `json:"content-version,omitempty"`
}

// Link represents a resource link
type Link struct {
	URL                 string `json:"URL,omitempty"`
	ContentType         string `json:"content-type,omitempty"`
	ContentVersion      string `json:"content-version,omitempty"`
	IntendedApplication string `json:"intended-application,omitempty"`
}

// Reference represents a reference
type Reference struct {
	Key                string `json:"key,omitempty"`
	DOI                string `json:"DOI,omitempty"`
	DOIAssertedBy      string `json:"doi-asserted-by,omitempty"`
	Issue              string `json:"issue,omitempty"`
	FirstPage          string `json:"first-page,omitempty"`
	Volume             string `json:"volume,omitempty"`
	Edition            string `json:"edition,omitempty"`
	Component          string `json:"component,omitempty"`
	StandardDesignator string `json:"standard-designator,omitempty"`
	StandardsBody      string `json:"standards-body,omitempty"`
	Author             string `json:"author,omitempty"`
	Year               string `json:"year,omitempty"`
	Unstructured       string `json:"unstructured,omitempty"`
	JournalTitle       string `json:"journal-title,omitempty"`
	ArticleTitle       string `json:"article-title,omitempty"`
	SeriesTitle        string `json:"series-title,omitempty"`
	VolumeTitle        string `json:"volume-title,omitempty"`
	ISSN               string `json:"ISSN,omitempty"`
	ISSNType           string `json:"issn-type,omitempty"`
	ISBN               string `json:"ISBN,omitempty"`
}

// Assertion represents an assertion
type Assertion struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
	URL   string `json:"URL,omitempty"`
	Label string `json:"label,omitempty"`
	Order int    `json:"order,omitempty"`
	Group *Group `json:"group,omitempty"`
}

// Group represents a group
type Group struct {
	Name  string `json:"name,omitempty"`
	Label string `json:"label,omitempty"`
}

// ContentDomain represents content domain
type ContentDomain struct {
	Domain               []string `json:"domain,omitempty"`
	CrossmarkRestriction bool     `json:"crossmark-restriction,omitempty"`
}

// ISSNType represents an ISSN type
type ISSNType struct {
	Value string `json:"value,omitempty"`
	Type  string `json:"type,omitempty"`
}

// JournalIssue represents a journal issue
type JournalIssue struct {
	Issue           string           `json:"issue,omitempty"`
	PublishedPrint  *PublishedPrint  `json:"published-print,omitempty"`
	PublishedOnline *PublishedOnline `json:"published-online,omitempty"`
}

// WorkResource represents a work resource
type WorkResource struct {
	Primary *Primary `json:"primary,omitempty"`
}

// Primary represents a primary resource
type Primary struct {
	URL string `json:"URL,omitempty"`
}

// DateParts represents a date with parts
type DateParts struct {
	DateParts [][]int   `json:"date-parts,omitempty"`
	DateTime  time.Time `json:"date-time"`
	Timestamp int64     `json:"timestamp,omitempty"`
}

// Indexed represents indexed date
type Indexed struct {
	DateParts [][]int   `json:"date-parts,omitempty"`
	DateTime  time.Time `json:"date-time"`
	Timestamp int64     `json:"timestamp,omitempty"`
}

// Created represents created date
type Created struct {
	DateParts [][]int   `json:"date-parts,omitempty"`
	DateTime  time.Time `json:"date-time"`
	Timestamp int64     `json:"timestamp,omitempty"`
}

// Deposited represents deposited date
type Deposited struct {
	DateParts [][]int   `json:"date-parts,omitempty"`
	DateTime  time.Time `json:"date-time"`
	Timestamp int64     `json:"timestamp,omitempty"`
}

// Issued represents issued date
type Issued struct {
	DateParts [][]int `json:"date-parts,omitempty"`
}

// Published represents published date
type Published struct {
	DateParts [][]int `json:"date-parts,omitempty"`
}

// PublishedPrint represents published print date
type PublishedPrint struct {
	DateParts [][]int `json:"date-parts,omitempty"`
}

// PublishedOnline represents published online date
type PublishedOnline struct {
	DateParts [][]int `json:"date-parts,omitempty"`
}

// WorkResponse wraps a single Work result
type WorkResponse struct {
	Status         string `json:"status"`
	MessageType    string `json:"message-type"`
	MessageVersion string `json:"message-version"`
	Message        Work   `json:"message"`
}

// MultipleWorksResponse wraps multiple Work results
type MultipleWorksResponse struct {
	Status         string               `json:"status"`
	MessageType    string               `json:"message-type"`
	MessageVersion string               `json:"message-version"`
	Message        MultipleWorksMessage `json:"message"`
}

// MultipleWorksMessage contains the list of works
type MultipleWorksMessage struct {
	Items        []Work  `json:"items"`
	TotalResults int     `json:"total-results"`
	ItemsPerPage int     `json:"items-per-page,omitempty"`
	Query        *Query  `json:"query,omitempty"`
	Facets       *Facets `json:"facets,omitempty"`
}

// JournalResponse wraps a single Journal result
type JournalResponse struct {
	Status         string  `json:"status"`
	MessageType    string  `json:"message-type"`
	MessageVersion string  `json:"message-version"`
	Message        Journal `json:"message"`
}

// MultipleJournalsResponse wraps multiple Journal results
type MultipleJournalsResponse struct {
	Status         string                  `json:"status"`
	MessageType    string                  `json:"message-type"`
	MessageVersion string                  `json:"message-version"`
	Message        MultipleJournalsMessage `json:"message"`
}

// MultipleJournalsMessage contains the list of journals
type MultipleJournalsMessage struct {
	Items        []Journal `json:"items"`
	TotalResults int       `json:"total-results"`
	ItemsPerPage int       `json:"items-per-page,omitempty"`
	Query        *Query    `json:"query,omitempty"`
}

// Query represents search query information
type Query struct {
	StartIndex  int `json:"start-index,omitempty"`
	SearchTerms any `json:"search-terms,omitempty"`
}

// Facets represents facets (empty placeholder for compatibility)
type Facets struct {
}
