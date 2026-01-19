package entities

import (
	"time"

	"github.com/google/uuid"
)

// RAiDInfo represents the RAiD metadata linked to a project.
type RAiDInfo struct {
	RAiDId                      string     `json:"raid_id"`
	SchemaUri                   string     `json:"schema_uri"`
	RegistrationAgencyId        string     `json:"registration_agency_id"`
	RegistrationAgencySchemaUri string     `json:"registration_agency_schema_uri"`
	OwnerId                     string     `json:"owner_id"`
	OwnerSchemaUri              string     `json:"owner_schema_uri"`
	OwnerServicePoint           *int64     `json:"owner_service_point,omitempty"`
	ProjectID                   uuid.UUID  `json:"project_id"`
	License                     string     `json:"license"`
	Version                     int        `json:"version"`
	LatestSync                  *time.Time `json:"latest_sync,omitempty"`
	Dirty                       bool       `json:"dirty"`
	Checksum                    *string    `json:"checksum,omitempty"`
}

// Default constants for RAiDInfo
const (
	RAiDSchemaUri = "https://raid.org/"
	RORSchemaUri  = "https://ror.org/"
	CC0License    = "Creative Commons CC-0"
)

// NewRAiDInfo creates a new RAiDInfo with default values.
func NewRAiDInfo(raidID, agencyID, ownerID string, projectID uuid.UUID) *RAiDInfo {
	return &RAiDInfo{
		RAiDId:                      raidID,
		SchemaUri:                   RAiDSchemaUri,
		RegistrationAgencyId:        agencyID,
		RegistrationAgencySchemaUri: RORSchemaUri,
		OwnerId:                     ownerID,
		OwnerSchemaUri:              RORSchemaUri,
		OwnerServicePoint:           nil,
		ProjectID:                   projectID,
		License:                     CC0License,
	}
}
