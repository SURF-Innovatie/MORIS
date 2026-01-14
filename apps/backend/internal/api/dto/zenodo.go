package dto

import ex "github.com/SURF-Innovatie/MORIS/external/zenodo"

type Deposition = ex.Deposition
type DepositionState = ex.DepositionState
type DepositionMetadata = ex.DepositionMetadata
type PrereserveDOI = ex.PrereserveDOI

type UploadType = ex.UploadType
type PublicationType = ex.PublicationType
type AccessRight = ex.AccessRight

type Creator = ex.Creator
type Contributor = ex.Contributor
type ContributorType = ex.ContributorType
type RelatedIdentifier = ex.RelatedIdentifier
type Community = ex.Community
type Grant = ex.Grant

type DepositionFile = ex.DepositionFile
type FileLinks = ex.FileLinks
type DepositionLinks = ex.DepositionLinks

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
		State:     d.State,
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
		UploadType:      m.UploadType,
		PublicationType: m.PublicationType,
		Description:     m.Description,
		AccessRight:     m.AccessRight,
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
				Type:        c.Type,
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
