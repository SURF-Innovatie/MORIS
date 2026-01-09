package raidsink

import (
	"context"
	"errors"

	"github.com/SURF-Innovatie/MORIS/external/raid"
	"github.com/SURF-Innovatie/MORIS/internal/adapter"
)

// RAiDSink exports projects to RAiD by minting new RAiDs.
type RAiDSink struct {
	client *raid.Client
	mapper *RAiDMapper
}

// NewRAiDSink creates a new RAiD sink adapter.
func NewRAiDSink(client *raid.Client, opts ...RAiDMapperOption) *RAiDSink {
	return &RAiDSink{
		client: client,
		mapper: NewRAiDMapper(opts...),
	}
}

func (s *RAiDSink) Name() string {
	return "raid"
}

func (s *RAiDSink) DisplayName() string {
	return "RAiD Registry"
}

func (s *RAiDSink) SupportedTypes() []adapter.DataType {
	return []adapter.DataType{adapter.DataTypeProject}
}

func (s *RAiDSink) OutputInfo() adapter.OutputInfo {
	return adapter.OutputInfo{
		Type:        adapter.TransferTypeAPI,
		Label:       "Mint RAiD",
		Description: "Register this project in the global RAiD registry",
	}
}

func (s *RAiDSink) Connect(ctx context.Context) error {
	// Client handles authentication, but we could check connectivity here
	return nil
}

// PushProject exports a project to RAiD by minting a new RAiD.
func (s *RAiDSink) PushProject(ctx context.Context, pc adapter.ProjectContext) error {
	if pc.Project == nil {
		return errors.New("project is required")
	}

	req := s.mapper.MapToCreateRequest(pc)
	_, err := s.client.MintRaid(ctx, req)
	return err
}

func (s *RAiDSink) PushUser(ctx context.Context, user adapter.UserContext) error {
	return errors.New("RAiD sink does not support user export")
}

// Close releases any resources.
func (s *RAiDSink) Close() error {
	return nil
}
