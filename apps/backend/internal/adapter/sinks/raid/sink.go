package raidsink

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/adapter"
)

// RAiDSink exports projects as RAiD-formatted JSON.
type RAiDSink struct {
	mapper *RAiDMapper
}

// NewRAiDSink creates a new RAiD sink adapter.
// The client parameter is ignored in this file-based version.
func NewRAiDSink(_ interface{}, opts ...RAiDMapperOption) *RAiDSink {
	return &RAiDSink{
		mapper: NewRAiDMapper(opts...),
	}
}

func (s *RAiDSink) Name() string {
	return "raid"
}

func (s *RAiDSink) DisplayName() string {
	return "RAiD Export (JSON)"
}

func (s *RAiDSink) SupportedTypes() []adapter.DataType {
	return []adapter.DataType{adapter.DataTypeProject}
}

func (s *RAiDSink) OutputInfo() adapter.OutputInfo {
	return adapter.OutputInfo{
		Type:        adapter.TransferTypeFile,
		Label:       "Download RAiD JSON",
		Description: "Download project as RAiD-formatted JSON file",
		MimeType:    "application/json",
	}
}

func (s *RAiDSink) Connect(ctx context.Context) error {
	return nil
}

// PushProject is not used for file-based exports.
func (s *RAiDSink) PushProject(ctx context.Context, pc adapter.ProjectContext) error {
	return errors.New("use ExportProjectData for file-based exports")
}

// ExportProjectData returns the project as RAiD-formatted JSON bytes.
func (s *RAiDSink) ExportProjectData(ctx context.Context, pc adapter.ProjectContext) (*adapter.ExportResult, error) {
	if pc.Project == nil {
		return nil, errors.New("project is required")
	}

	req := s.mapper.MapToCreateRequest(pc)

	data, err := json.MarshalIndent(req, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	filename := fmt.Sprintf("raid_%s_%s.json",
		pc.ProjectID.String(),
		time.Now().Format("20060102_150405"),
	)

	return &adapter.ExportResult{
		Data:     data,
		Filename: filename,
		MimeType: "application/json",
	}, nil
}

func (s *RAiDSink) PushUser(ctx context.Context, user adapter.UserContext) error {
	return errors.New("RAiD sink does not support user export")
}

func (s *RAiDSink) Close() error {
	return nil
}
