package dto

type OrcidPerson struct {
	FirstName  string `json:"first_name,omitempty"`
	LastName   string `json:"last_name,omitempty"`
	CreditName string `json:"credit_name,omitempty"`
	Biography  string `json:"biography,omitempty"`
	ORCID      string `json:"orcid,omitempty"`
}
