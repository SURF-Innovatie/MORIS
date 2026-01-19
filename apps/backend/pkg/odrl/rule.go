package odrl

import "encoding/json"

// AssetRef can hold either a simple IRI string or a full Asset object.
// When marshaling, it outputs either a string or an object depending on which is set.
type AssetRef struct {
	IRI   string // Simple IRI reference
	Asset *Asset // Full Asset object
}

func (a AssetRef) MarshalJSON() ([]byte, error) {
	if a.Asset != nil {
		return json.Marshal(a.Asset)
	}
	return json.Marshal(a.IRI)
}

func (a *AssetRef) UnmarshalJSON(data []byte) error {
	// Try string first
	var iri string
	if err := json.Unmarshal(data, &iri); err == nil {
		a.IRI = iri
		return nil
	}
	// Try object
	var asset Asset
	if err := json.Unmarshal(data, &asset); err == nil {
		a.Asset = &asset
		return nil
	}
	return nil
}

// AssetIRI creates an AssetRef from a simple IRI string.
func AssetIRI(iri string) AssetRef {
	return AssetRef{IRI: iri}
}

// AssetObject creates an AssetRef from a full Asset object.
func AssetObject(asset *Asset) AssetRef {
	return AssetRef{Asset: asset}
}

// PartyRef can hold either a simple IRI string or a full Party object.
type PartyRef struct {
	IRI   string
	Party *Party
}

func (p PartyRef) MarshalJSON() ([]byte, error) {
	if p.Party != nil {
		return json.Marshal(p.Party)
	}
	return json.Marshal(p.IRI)
}

func (p *PartyRef) UnmarshalJSON(data []byte) error {
	var iri string
	if err := json.Unmarshal(data, &iri); err == nil {
		p.IRI = iri
		return nil
	}
	var party Party
	if err := json.Unmarshal(data, &party); err == nil {
		p.Party = &party
		return nil
	}
	return nil
}

// PartyIRI creates a PartyRef from a simple IRI string.
func PartyIRI(iri string) PartyRef {
	return PartyRef{IRI: iri}
}

// PartyObject creates a PartyRef from a full Party object.
func PartyObject(party *Party) PartyRef {
	return PartyRef{Party: party}
}

// ActionRef can hold either a simple action name/IRI or a refined Action object.
type ActionRef struct {
	Name   string  // Simple action name (e.g., "use", "play")
	Action *Action // Refined action with constraints
}

func (a ActionRef) MarshalJSON() ([]byte, error) {
	if a.Action != nil {
		return json.Marshal(a.Action)
	}
	return json.Marshal(a.Name)
}

func (a *ActionRef) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err == nil {
		a.Name = name
		return nil
	}
	var action Action
	if err := json.Unmarshal(data, &action); err == nil {
		a.Action = &action
		return nil
	}
	return nil
}

// ActionName creates an ActionRef from a simple action name.
func ActionName(name string) ActionRef {
	return ActionRef{Name: name}
}

// ActionRefined creates an ActionRef from a refined Action object.
func ActionRefined(action *Action) ActionRef {
	return ActionRef{Action: action}
}

// RuleCommon holds properties common to all Rules.
type RuleCommon struct {
	UID        string        `json:"uid,omitempty"`
	Action     ActionRef     `json:"action"`
	Target     *AssetRef     `json:"target,omitempty"`
	Assigner   *PartyRef     `json:"assigner,omitempty"`
	Assignee   *PartyRef     `json:"assignee,omitempty"`
	Constraint []*Constraint `json:"constraint,omitempty"`
}

// Permission allows an action to be exercised on an Asset.
type Permission struct {
	RuleCommon
	Duty []*Duty `json:"duty,omitempty"`
}

// Prohibition disallows an action.
type Prohibition struct {
	RuleCommon
}

// Duty represents an obligation to exercise an action.
type Duty struct {
	RuleCommon
	Consequence []*Duty `json:"consequence,omitempty"`
}

// NewPermission creates a new Permission with a simple action name.
func NewPermission(action string) *Permission {
	return &Permission{
		RuleCommon: RuleCommon{
			Action: ActionName(action),
		},
	}
}

// NewPermissionWithAction creates a Permission with a refined Action.
func NewPermissionWithAction(action *Action) *Permission {
	return &Permission{
		RuleCommon: RuleCommon{
			Action: ActionRefined(action),
		},
	}
}

// NewProhibition creates a new Prohibition with a simple action name.
func NewProhibition(action string) *Prohibition {
	return &Prohibition{
		RuleCommon: RuleCommon{
			Action: ActionName(action),
		},
	}
}

// NewDuty creates a new Duty with a simple action name.
func NewDuty(action string) *Duty {
	return &Duty{
		RuleCommon: RuleCommon{
			Action: ActionName(action),
		},
	}
}

// Permission Builders

// WithTarget sets the target asset using a simple IRI.
func (p *Permission) WithTarget(iri string) *Permission {
	p.Target = &AssetRef{IRI: iri}
	return p
}

// WithTargetAsset sets the target using a full Asset object.
func (p *Permission) WithTargetAsset(asset *Asset) *Permission {
	p.Target = &AssetRef{Asset: asset}
	return p
}

// WithAssigner sets the assigner using a simple IRI.
func (p *Permission) WithAssigner(iri string) *Permission {
	p.Assigner = &PartyRef{IRI: iri}
	return p
}

// WithAssignerParty sets the assigner using a full Party object.
func (p *Permission) WithAssignerParty(party *Party) *Permission {
	p.Assigner = &PartyRef{Party: party}
	return p
}

// WithAssignee sets the assignee using a simple IRI.
func (p *Permission) WithAssignee(iri string) *Permission {
	p.Assignee = &PartyRef{IRI: iri}
	return p
}

// WithAssigneeParty sets the assignee using a full Party object.
func (p *Permission) WithAssigneeParty(party *Party) *Permission {
	p.Assignee = &PartyRef{Party: party}
	return p
}

// WithConstraint adds a constraint to the permission.
func (p *Permission) WithConstraint(c *Constraint) *Permission {
	p.Constraint = append(p.Constraint, c)
	return p
}

// WithDuty adds a duty to the permission.
func (p *Permission) WithDuty(d *Duty) *Permission {
	p.Duty = append(p.Duty, d)
	return p
}

// Prohibition Builders

func (p *Prohibition) WithTarget(iri string) *Prohibition {
	p.Target = &AssetRef{IRI: iri}
	return p
}

func (p *Prohibition) WithTargetAsset(asset *Asset) *Prohibition {
	p.Target = &AssetRef{Asset: asset}
	return p
}

func (p *Prohibition) WithAssigner(iri string) *Prohibition {
	p.Assigner = &PartyRef{IRI: iri}
	return p
}

func (p *Prohibition) WithAssignee(iri string) *Prohibition {
	p.Assignee = &PartyRef{IRI: iri}
	return p
}

func (p *Prohibition) WithConstraint(c *Constraint) *Prohibition {
	p.Constraint = append(p.Constraint, c)
	return p
}

// Duty Builders

func (d *Duty) WithTarget(iri string) *Duty {
	d.Target = &AssetRef{IRI: iri}
	return d
}

func (d *Duty) WithAssigner(iri string) *Duty {
	d.Assigner = &PartyRef{IRI: iri}
	return d
}

func (d *Duty) WithAssignee(iri string) *Duty {
	d.Assignee = &PartyRef{IRI: iri}
	return d
}

func (d *Duty) WithConstraint(c *Constraint) *Duty {
	d.Constraint = append(d.Constraint, c)
	return d
}

func (d *Duty) WithConsequence(cons *Duty) *Duty {
	d.Consequence = append(d.Consequence, cons)
	return d
}
