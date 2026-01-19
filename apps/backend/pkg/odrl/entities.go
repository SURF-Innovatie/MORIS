package odrl

// Asset represents a resource.
type Asset struct {
	UID    string   `json:"uid,omitempty"`
	Type   string   `json:"@type,omitempty"` // Typically "Asset" or "AssetCollection"
	PartOf []string `json:"partOf,omitempty"`
}

// NewAsset creates a new Asset.
func NewAsset(uid string) *Asset {
	return &Asset{
		UID:  uid,
		Type: "Asset",
	}
}

// NewAssetCollection creates a new Asset that is a collection.
func NewAssetCollection(uid string) *Asset {
	return &Asset{
		UID:  uid,
		Type: "AssetCollection",
	}
}

// Party represents an entity.
type Party struct {
	UID    string   `json:"uid,omitempty"`
	Type   string   `json:"@type,omitempty"` // Typically "Party" or "PartyCollection"
	PartOf []string `json:"partOf,omitempty"`
}

// NewParty creates a new Party.
func NewParty(uid string) *Party {
	return &Party{
		UID:  uid,
		Type: "Party",
	}
}

// NewPartyCollection creates a new Party that is a collection.
func NewPartyCollection(uid string) *Party {
	return &Party{
		UID:  uid,
		Type: "PartyCollection",
	}
}

// Action represents an operation.
// In simple cases, it's just a string (IRI).
// In complex cases (refined action), it can be a struct.
type Action struct {
	Value      string        `json:"rdf:value,omitempty"` // If serialization needs value key? Actually ODRL often uses just string for action.
	Refinement []*Constraint `json:"refinement,omitempty"`
	IncludedIn string        `json:"includedIn,omitempty"`
	Implies    string        `json:"implies,omitempty"`
}

// Custom Marshaling for Action might be needed if we want to support both string and object.
// But for now, let's keep it simple. The Rule struct uses interface{} for Action.

// Common ODRL Actions
const (
	ActionUse          = "use"
	ActionTransfer     = "transfer"
	ActionAttribution  = "attribution"
	ActionCompensate   = "compensate"
	ActionDelete       = "delete"
	ActionDistribute   = "distribute"
	ActionExecute      = "execute"
	ActionExtract      = "extract"
	ActionGrantUse     = "grantUse"
	ActionIndex        = "index"
	ActionInstall      = "install"
	ActionModify       = "modify"
	ActionMove         = "move"
	ActionPlay         = "play"
	ActionPresent      = "present"
	ActionPrint        = "print"
	ActionRead         = "read"
	ActionReproduce    = "reproduce"
	ActionReviewPolicy = "reviewPolicy"
	ActionSell         = "sell"
	ActionStream       = "stream"
	ActionSynchronize  = "synchronize"
	ActionTextToSpeech = "textToSpeech"
	ActionTransform    = "transform"
	ActionTranslate    = "translate"
)
