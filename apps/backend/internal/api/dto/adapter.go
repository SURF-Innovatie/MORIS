package dto

import "github.com/SURF-Innovatie/MORIS/internal/adapter"

type AdapterInfoResponse struct {
	Name           string             `json:"name"`
	DisplayName    string             `json:"display_name"`
	SupportedTypes []adapter.DataType `json:"supported_types"`
}

type SourceInfoResponse struct {
	AdapterInfoResponse
	Input adapter.InputInfo `json:"input"`
}

type SinkInfoResponse struct {
	AdapterInfoResponse
	Output adapter.OutputInfo `json:"output"`
}

type AdapterListResponse struct {
	Sources []SourceInfoResponse `json:"sources"`
	Sinks   []SinkInfoResponse   `json:"sinks"`
}
