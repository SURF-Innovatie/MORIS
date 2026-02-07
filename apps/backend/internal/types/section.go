package types

// Section represents a content section in a Page.
// It matches the JSON structure: { "type": "...", "data": ... }
type Section struct {
	ID   string                 `json:"id"`
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}
