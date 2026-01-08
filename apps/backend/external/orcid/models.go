package orcid

// OrcidPerson represents a person record from ORCID
type OrcidPerson struct {
	FirstName  string `json:"first_name,omitempty"`
	LastName   string `json:"last_name,omitempty"`
	CreditName string `json:"credit_name,omitempty"`
	Biography  string `json:"biography,omitempty"`
	ORCID      string `json:"orcid,omitempty"`
}

// PersonExpandedSearchResult represents an item in the expanded search results
type PersonExpandedSearchResult struct {
	OrcidID    string `json:"orcid-id"`
	GivenNames string `json:"given-names"`
	FamilyNames string `json:"family-names"`
	CreditName string `json:"credit-name"`
}

// ToPerson converts the search result to a domain person object
func (p *PersonExpandedSearchResult) ToPerson() OrcidPerson {
	return OrcidPerson{
		FirstName:  p.GivenNames,
		LastName:   p.FamilyNames,
		CreditName: p.CreditName,
		ORCID:      p.OrcidID,
	}
}

// expandedSearchResponse represents the top-level response for an expanded search
type expandedSearchResponse struct {
	ExpandedResult []PersonExpandedSearchResult `json:"expanded-result"`
}
