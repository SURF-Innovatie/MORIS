package dto

type AvailableEvent struct {
	Type          string         `json:"type"`
	FriendlyName  string         `json:"friendly_name"`
	NeedsApproval bool           `json:"needs_approval"`
	Allowed       bool           `json:"allowed"`
	InputSchema   map[string]any `json:"input_schema,omitempty"` // optional if you add schemas
}
