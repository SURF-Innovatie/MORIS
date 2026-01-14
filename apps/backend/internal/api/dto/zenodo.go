package dto

import "time"
import ex "github.com/SURF-Innovatie/MORIS/external/zenodo"

// Deposition represents a Zenodo deposition resource.
type Deposition struct {
	ID        int                 `json:"id"`
	ConceptID string              `json:"conceptrecid,omitempty"`
	DOI       string              `json:"doi,omitempty"`
	DOIURL    string              `json:"doi_url,omitempty"`
	RecordID  int                 `json:"record_id,omitempty"`
	RecordURL string              `json:"record_url,omitempty"`
	Created   time.Time           `json:"created,omitempty"`
	Modified  time.Time           `json:"modified,omitempty"`
	Owner     int                 `json:"owner,omitempty"`
	State     DepositionState     `json:"state,omitempty"`
	Submitted bool                `json:"submitted,omitempty"`
	Title     string              `json:"title,omitempty"`
	Metadata  *DepositionMetadata `json:"metadata,omitempty"`
	Files     []DepositionFile    `json:"files,omitempty"`
	Links     *DepositionLinks    `json:"links,omitempty"`
}

type DepositionState string

const (
	StateInProgress  DepositionState = "inprogress"
	StateDone        DepositionState = "done"
	StateError       DepositionState = "error"
	StateUnsubmitted DepositionState = "unsubmitted"
)

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

type PrereserveDOI struct {
	DOI   string `json:"doi"`
	RecID int    `json:"recid"`
}

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

type AccessRight string

const (
	AccessOpen       AccessRight = "open"
	AccessEmbargoed  AccessRight = "embargoed"
	AccessRestricted AccessRight = "restricted"
	AccessClosed     AccessRight = "closed"
)

type Creator struct {
	Name        string `json:"name"`
	Affiliation string `json:"affiliation,omitempty"`
	ORCID       string `json:"orcid,omitempty"`
	GND         string `json:"gnd,omitempty"`
}

type Contributor struct {
	Name        string          `json:"name"`
	Type        ContributorType `json:"type"`
	Affiliation string          `json:"affiliation,omitempty"`
	ORCID       string          `json:"orcid,omitempty"`
	GND         string          `json:"gnd,omitempty"`
}

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

type RelatedIdentifier struct {
	Identifier   string `json:"identifier"`
	Relation     string `json:"relation"`
	ResourceType string `json:"resource_type,omitempty"`
}

type Community struct {
	Identifier string `json:"identifier"`
}

type Grant struct {
	ID string `json:"id"`
}

type DepositionFile struct {
	ID       string     `json:"id"`
	Filename string     `json:"filename"`
	Filesize int64      `json:"filesize"`
	Checksum string     `json:"checksum"`
	Links    *FileLinks `json:"links,omitempty"`
}

type FileLinks struct {
	Self     string `json:"self,omitempty"`
	Download string `json:"download,omitempty"`
}

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

func FromExternalDeposition(d ex.Deposition) Deposition {
	out := Deposition{
		ID:        d.ID,
		ConceptID: d.ConceptID,
		DOI:       d.DOI,
		DOIURL:    d.DOIURL,
		RecordID:  d.RecordID,
		RecordURL: d.RecordURL,
		Created:   d.Created,
		Modified:  d.Modified,
		Owner:     d.Owner,
		State:     DepositionState(d.State),
		Submitted: d.Submitted,
		Title:     d.Title,
	}

	if d.Metadata != nil {
		out.Metadata = FromExternalDepositionMetadata(*d.Metadata)
	}
	if d.Links != nil {
		out.Links = FromExternalDepositionLinks(*d.Links)
	}
	if len(d.Files) > 0 {
		out.Files = make([]DepositionFile, 0, len(d.Files))
		for _, f := range d.Files {
			out.Files = append(out.Files, FromExternalDepositionFile(f))
		}
	}

	return out
}

func FromExternalDepositions(ds []ex.Deposition) []Deposition {
	out := make([]Deposition, 0, len(ds))
	for _, d := range ds {
		out = append(out, FromExternalDeposition(d))
	}
	return out
}

func FromExternalDepositionMetadata(m ex.DepositionMetadata) *DepositionMetadata {
	out := DepositionMetadata{
		Title:           m.Title,
		UploadType:      UploadType(m.UploadType),
		PublicationType: PublicationType(m.PublicationType),
		Description:     m.Description,
		AccessRight:     AccessRight(m.AccessRight),
		License:         m.License,
		DOI:             m.DOI,
		Keywords:        m.Keywords,
		Notes:           m.Notes,
		References:      m.References,
		PublicationDate: m.PublicationDate,
	}

	if m.PrereserveDOI != nil {
		out.PrereserveDOI = &PrereserveDOI{DOI: m.PrereserveDOI.DOI, RecID: m.PrereserveDOI.RecID}
	}

	if len(m.Creators) > 0 {
		out.Creators = make([]Creator, 0, len(m.Creators))
		for _, c := range m.Creators {
			out.Creators = append(out.Creators, Creator{
				Name:        c.Name,
				Affiliation: c.Affiliation,
				ORCID:       c.ORCID,
				GND:         c.GND,
			})
		}
	}

	if len(m.Contributors) > 0 {
		out.Contributors = make([]Contributor, 0, len(m.Contributors))
		for _, c := range m.Contributors {
			out.Contributors = append(out.Contributors, Contributor{
				Name:        c.Name,
				Type:        ContributorType(c.Type),
				Affiliation: c.Affiliation,
				ORCID:       c.ORCID,
				GND:         c.GND,
			})
		}
	}

	if len(m.RelatedIDs) > 0 {
		out.RelatedIDs = make([]RelatedIdentifier, 0, len(m.RelatedIDs))
		for _, r := range m.RelatedIDs {
			out.RelatedIDs = append(out.RelatedIDs, RelatedIdentifier{
				Identifier:   r.Identifier,
				Relation:     r.Relation,
				ResourceType: r.ResourceType,
			})
		}
	}

	if len(m.Communities) > 0 {
		out.Communities = make([]Community, 0, len(m.Communities))
		for _, c := range m.Communities {
			out.Communities = append(out.Communities, Community{Identifier: c.Identifier})
		}
	}

	if len(m.Grants) > 0 {
		out.Grants = make([]Grant, 0, len(m.Grants))
		for _, g := range m.Grants {
			out.Grants = append(out.Grants, Grant{ID: g.ID})
		}
	}

	return &out
}

func FromExternalDepositionLinks(l ex.DepositionLinks) *DepositionLinks {
	return &DepositionLinks{
		Bucket:          l.Bucket,
		Discard:         l.Discard,
		Edit:            l.Edit,
		Files:           l.Files,
		HTML:            l.HTML,
		LatestDraft:     l.LatestDraft,
		LatestDraftHTML: l.LatestDraftHTML,
		Publish:         l.Publish,
		Self:            l.Self,
		NewVersion:      l.NewVersion,
	}
}

func FromExternalDepositionFile(f ex.DepositionFile) DepositionFile {
	out := DepositionFile{
		ID:       f.ID,
		Filename: f.Filename,
		Filesize: f.Filesize,
		Checksum: f.Checksum,
	}
	if f.Links != nil {
		out.Links = &FileLinks{Self: f.Links.Self, Download: f.Links.Download}
	}
	return out
}

func FromExternalDepositionFilePtr(f *ex.DepositionFile) *DepositionFile {
	if f == nil {
		return nil
	}
	ff := FromExternalDepositionFile(*f)
	return &ff
}
