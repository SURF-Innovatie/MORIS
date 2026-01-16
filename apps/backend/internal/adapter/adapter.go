package adapter

import (
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/google/uuid"
)

// DataType represents the type of data being transferred.
type DataType string

const (
	DataTypeUser    DataType = "user"
	DataTypeProject DataType = "project"
)

// TransferType indicates how the data is delivered/received (API or File)
type TransferType string

const (
	TransferTypeAPI  TransferType = "api"  // External service call
	TransferTypeFile TransferType = "file" // File upload/download
)

// InputInfo describes a source adapter's input requirements
type InputInfo struct {
	Type        TransferType `json:"type"`
	Label       string       `json:"label"`
	Description string       `json:"description"`
	MimeType    string       `json:"mime_type,omitempty"` // For file inputs
}

// OutputInfo describes a sink adapter's output behavior
type OutputInfo struct {
	Type        TransferType `json:"type"`
	Label       string       `json:"label"`
	Description string       `json:"description"`
	MimeType    string       `json:"mime_type,omitempty"` // For file outputs
}

// ProjectContext bundles a project's event stream with enriched entity data.
// This is what sinks consume - the full event history plus resolved entities.
type ProjectContext struct {
	ProjectID uuid.UUID
	Events    []events.Event             // The raw event stream
	Project   *entities.Project          // Reduced/projected state
	Members   []entities.Person          // Resolved member entities
	OrgNode   *entities.OrganisationNode // Owning organisation
}

// UserContext bundles user data for import/export.
type UserContext struct {
	Person *entities.Person
	User   *entities.User
}

// FetchOptions configures how data is fetched from a source
type FetchOptions struct {
	BatchSize int                    // Number of records per batch
	Filters   map[string]interface{} // Optional filters
}
