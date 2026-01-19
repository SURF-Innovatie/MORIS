package odrl

import "encoding/json"

const (
	// ODRLContextURL is the standard ODRL JSON-LD context URL.
	ODRLContextURL = "http://www.w3.org/ns/odrl.jsonld"
)

// PolicyType represents the subclass of Policy (Set, Offer, Agreement).
type PolicyType string

const (
	PolicyTypeSet       PolicyType = "Set"
	PolicyTypeOffer     PolicyType = "Offer"
	PolicyTypeAgreement PolicyType = "Agreement"
	PolicyTypePrivacy   PolicyType = "Privacy" // Valid for ODRL? Not standard subclass but widely used. Stick to standard 3 for now.
)

// Policy represents the ODRL Policy class.
type Policy struct {
	Context     string         `json:"@context,omitempty"`
	Type        PolicyType     `json:"@type"`
	UID         string         `json:"uid,omitempty"`
	Profile     string         `json:"profile,omitempty"`
	InheritFrom string         `json:"inheritFrom,omitempty"`
	Permission  []*Permission  `json:"permission,omitempty"`
	Prohibition []*Prohibition `json:"prohibition,omitempty"`
	Obligation  []*Duty        `json:"obligation,omitempty"` // "obligation" is the property name for Duty in Policy
}

// NewPolicy creates a new Policy with the given type and UID.
// It automatically sets the Context.
func NewPolicy(t PolicyType, uid string) *Policy {
	return &Policy{
		Context: ODRLContextURL,
		Type:    t,
		UID:     uid,
	}
}

// NewSet creates a new Set Policy.
func NewSet(uid string) *Policy {
	return NewPolicy(PolicyTypeSet, uid)
}

// NewOffer creates a new Offer Policy.
func NewOffer(uid string) *Policy {
	return NewPolicy(PolicyTypeOffer, uid)
}

// NewAgreement creates a new Agreement Policy.
func NewAgreement(uid string) *Policy {
	return NewPolicy(PolicyTypeAgreement, uid)
}

// WithProfile sets the profile property.
func (p *Policy) WithProfile(profile string) *Policy {
	p.Profile = profile
	return p
}

// WithInheritFrom sets the inheritFrom property.
func (p *Policy) WithInheritFrom(inheritFrom string) *Policy {
	p.InheritFrom = inheritFrom
	return p
}

// AddPermission adds a permission rule to the policy.
func (p *Policy) AddPermission(perm *Permission) *Policy {
	p.Permission = append(p.Permission, perm)
	return p
}

// AddProhibition adds a prohibition rule to the policy.
func (p *Policy) AddProhibition(prohib *Prohibition) *Policy {
	p.Prohibition = append(p.Prohibition, prohib)
	return p
}

// AddObligation adds a duty rule to the policy.
func (p *Policy) AddObligation(duty *Duty) *Policy {
	p.Obligation = append(p.Obligation, duty)
	return p
}

// ToJSON serializes the Policy to JSON-LD.
func (p *Policy) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}

// ToJSONIndent serializes the Policy to indented JSON-LD.
func (p *Policy) ToJSONIndent() ([]byte, error) {
	return json.MarshalIndent(p, "", "  ")
}
