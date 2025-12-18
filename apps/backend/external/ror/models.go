package ror

// OrganizationStatus represents the status of an organization.
type OrganizationStatus string

const (
	OrganizationStatusActive    OrganizationStatus = "active"
	OrganizationStatusInactive  OrganizationStatus = "inactive"
	OrganizationStatusWithdrawn OrganizationStatus = "withdrawn"
)

// OrganizationType represents the type of an organization.
type OrganizationType string

const (
	OrganizationTypeEducation  OrganizationType = "education"
	OrganizationTypeFunder     OrganizationType = "funder"
	OrganizationTypeHealthcare OrganizationType = "healthcare"
	OrganizationTypeCompany    OrganizationType = "company"
	OrganizationTypeArchive    OrganizationType = "archive"
	OrganizationTypeNonprofit  OrganizationType = "nonprofit"
	OrganizationTypeGovernment OrganizationType = "government"
	OrganizationTypeFacility   OrganizationType = "facility"
	OrganizationTypeOther      OrganizationType = "other"
)

// OrganizationRelationshipType represents the type of relationship.
type OrganizationRelationshipType string

const (
	OrganizationRelationshipTypeRelated     OrganizationRelationshipType = "Related"
	OrganizationRelationshipTypeParent      OrganizationRelationshipType = "Parent"
	OrganizationRelationshipTypeChild       OrganizationRelationshipType = "Child"
	OrganizationRelationshipTypePredecessor OrganizationRelationshipType = "Predecessor"
	OrganizationRelationshipTypeSuccessor   OrganizationRelationshipType = "Successor"
)

// OrganizationNameType represents the type of organization name.
type OrganizationNameType string

const (
	OrganizationNameTypeAcronym    OrganizationNameType = "acronym"
	OrganizationNameTypeAlias      OrganizationNameType = "alias"
	OrganizationNameTypeLabel      OrganizationNameType = "label"
	OrganizationNameTypeRorDisplay OrganizationNameType = "ror_display"
)

// OrganizationName represents an alternative name for an organization.
type OrganizationName struct {
	Value string                 `json:"value"`
	Types []OrganizationNameType `json:"types"`
	Lang  *string                `json:"lang,omitempty"`
}

// OrganizationRelationship represents a relationship to another organization.
type OrganizationRelationship struct {
	Id    string                       `json:"id"`
	Label string                       `json:"label"`
	Type  OrganizationRelationshipType `json:"type"`
}

// GeonamesDetails represents location details from Geonames.
type GeonamesDetails struct {
	CountryCode            *string  `json:"country_code,omitempty"`
	CountryName            *string  `json:"country_name,omitempty"`
	CountrySubdivisionCode *string  `json:"country_subdivision_code,omitempty"`
	CountrySubdivisionName *string  `json:"country_subdivision_name,omitempty"`
	Latitude               *float64 `json:"lat,omitempty"`
	Longitude              *float64 `json:"lng,omitempty"`
	Name                   string   `json:"name"`
}

// OrganizationLocation represents a location of an organization.
type OrganizationLocation struct {
	GeonamesId      int             `json:"geonames_id"`
	GeonamesDetails GeonamesDetails `json:"geonames_details"`
}

// Organization represents a Research Organization Registry (ROR) organization.
type Organization struct {
	Id                 string                     `json:"id"`
	Name               string                     `json:"name"`
	Domains            []string                   `json:"domains,omitempty"`
	Established        *int                       `json:"established,omitempty"`
	Links              []string                   `json:"links,omitempty"`
	Relationships      []OrganizationRelationship `json:"relationships,omitempty"`
	OrganizationStatus OrganizationStatus         `json:"status"`
	Types              []OrganizationType         `json:"types"`
	// Field Locations is not present in the reference C# Organization.cs file,
	// but OrganizationLocation.cs was provided and the QueryBuilder filters on locations.
	Locations []OrganizationLocation `json:"locations,omitempty"`
}

// MetadataCounter is an interface for types that have an ID and a generic count.
type MetadataCounter interface {
	GetId() string
	GetCount() int
	AddCount(n int)
}

// MetadataCount represents a generic count in metadata.
type MetadataCount struct {
	Count int    `json:"count"`
	Id    string `json:"id"`
	Title string `json:"title"`
}

func (m *MetadataCount) GetId() string  { return m.Id }
func (m *MetadataCount) GetCount() int  { return m.Count }
func (m *MetadataCount) AddCount(n int) { m.Count += n }

// MetadataContinentCount represents a count of organizations in a continent.
type MetadataContinentCount struct {
	MetadataCount
}

// MetadataCountryCount represents a count of organizations in a country.
type MetadataCountryCount struct {
	MetadataCount
}

// MetadataStatusCount represents a count of organizations with a specific status.
type MetadataStatusCount struct {
	MetadataCount
}

// MetadataTypeCount represents a count of organizations with a specific type.
type MetadataTypeCount struct {
	MetadataCount
}

// ResultMetadata represents metadata about the query results.
type ResultMetadata struct {
	Continents []MetadataContinentCount `json:"continents"`
	Countries  []MetadataCountryCount   `json:"countries"`
	Statuses   []MetadataStatusCount    `json:"statuses"`
	Types      []MetadataTypeCount      `json:"types"`
}

// OrganizationsResult represents the result of an organization query.
type OrganizationsResult struct {
	Organizations   []Organization `json:"items"`
	Metadata        ResultMetadata `json:"meta"`
	NumberOfResults int            `json:"number_of_results"`
	TimeTaken       int            `json:"time_taken"`
}

// Combine merges two OrganizationsResult objects.
func (r *OrganizationsResult) Combine(other *OrganizationsResult) *OrganizationsResult {
	if r == nil {
		return other
	}
	if other == nil {
		return r
	}

	combined := &OrganizationsResult{
		Organizations:   append(r.Organizations, other.Organizations...),
		NumberOfResults: r.NumberOfResults + other.NumberOfResults,
		TimeTaken:       r.TimeTaken + other.TimeTaken,
		Metadata: ResultMetadata{
			Continents: collectCounts(r.Metadata.Continents, other.Metadata.Continents),
			Countries:  collectCounts(r.Metadata.Countries, other.Metadata.Countries),
			Statuses:   collectCounts(r.Metadata.Statuses, other.Metadata.Statuses),
			Types:      collectCounts(r.Metadata.Types, other.Metadata.Types),
		},
	}
	return combined
}

func collectCounts[T any, PT interface {
	*T
	MetadataCounter
}](c1 []T, c2 []T) []T {
	counts := make(map[string]T)

	for _, item := range c1 {
		p := PT(&item)
		counts[p.GetId()] = item
	}

	for _, item := range c2 {
		p := PT(&item)
		id := p.GetId()
		if existing, ok := counts[id]; ok {
			ePtr := PT(&existing)
			ePtr.AddCount(p.GetCount())
			counts[id] = existing
		} else {
			counts[id] = item
		}
	}

	result := make([]T, 0, len(counts))
	for _, v := range counts {
		result = append(result, v)
	}
	return result
}
