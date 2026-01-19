package dto

import ex "github.com/SURF-Innovatie/MORIS/external/raid"

type RAiD = ex.RAiDDto
type RAiDCreateRequest = ex.RAiDCreateRequest
type RAiDUpdateRequest = ex.RAiDUpdateRequest

// Helper alias types for sub-structs to ensure they are available in the API package if needed
type RAiDId = ex.RAiDId
type RAiDOwner = ex.RAiDOwner
type RAiDTitle = ex.RAiDTitle
type RAiDDate = ex.RAiDDate
type RAiDDescription = ex.RAiDDescription
type RAiDAccess = ex.RAiDAccess
type RAiDAlternateUrl = ex.RAiDAlternateUrl
type RAiDContributor = ex.RAiDContributor
type RAiDOrganisation = ex.RAiDOrganisation
type RAiDSubject = ex.RAiDSubject
type RAiDRelatedRaid = ex.RAiDRelatedRaid
type RAiDRelatedObject = ex.RAiDRelatedObject
type RAiDAlternateIdentifier = ex.RAiDAlternateIdentifier
type RAiDSpatialCoverage = ex.RAiDSpatialCoverage
type RAiDTraditionalKnowledgeLabel = ex.RAiDTraditionalKnowledgeLabel
type RAiDMetadata = ex.RAiDMetadata

func FromExternalRaid(r ex.RAiDDto) RAiD {
	// Since we are aliasing directly, we can just cast or return.
	// However, if we need deep copy or specific transformation later, we can add it here.
	// For now, direct casting/assignment is sufficient as types are identical.
	return r
}

func FromExternalRaids(rs []ex.RAiDDto) []RAiD {
	out := make([]RAiD, 0, len(rs))
	for _, r := range rs {
		out = append(out, FromExternalRaid(r))
	}
	return out
}
