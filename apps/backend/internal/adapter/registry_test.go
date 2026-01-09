package adapter_test

import (
	"context"
	"testing"

	"github.com/SURF-Innovatie/MORIS/internal/adapter"
)

// MockSinkAdapter for testing
type mockSinkAdapter struct {
	name string
}

func (m *mockSinkAdapter) Name() string                                                       { return m.name }
func (m *mockSinkAdapter) DisplayName() string                                                { return m.name }
func (m *mockSinkAdapter) SupportedTypes() []adapter.DataType                                 { return []adapter.DataType{adapter.DataTypeProject} }
func (m *mockSinkAdapter) OutputInfo() adapter.OutputInfo                                     { return adapter.OutputInfo{Type: adapter.TransferTypeAPI} }
func (m *mockSinkAdapter) Connect(ctx context.Context) error                                  { return nil }
func (m *mockSinkAdapter) PushProject(ctx context.Context, pc adapter.ProjectContext) error   { return nil }
func (m *mockSinkAdapter) ExportProjectData(ctx context.Context, pc adapter.ProjectContext) (*adapter.ExportResult, error) { return nil, nil }
func (m *mockSinkAdapter) PushUser(ctx context.Context, uc adapter.UserContext) error         { return nil }
func (m *mockSinkAdapter) Close() error                                                       { return nil }

func TestRegistry_RegisterSink(t *testing.T) {
	r := adapter.NewRegistry()

	sink := &mockSinkAdapter{name: "test-sink"}
	err := r.RegisterSink(sink)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify we can retrieve it
	retrieved, ok := r.GetSink("test-sink")
	if !ok {
		t.Fatal("expected to find registered sink")
	}
	if retrieved.Name() != "test-sink" {
		t.Errorf("expected name 'test-sink', got %q", retrieved.Name())
	}
}

func TestRegistry_RegisterSink_Duplicate(t *testing.T) {
	r := adapter.NewRegistry()

	sink1 := &mockSinkAdapter{name: "dup-sink"}
	sink2 := &mockSinkAdapter{name: "dup-sink"}

	_ = r.RegisterSink(sink1)
	err := r.RegisterSink(sink2)
	if err == nil {
		t.Fatal("expected error for duplicate registration")
	}
}

func TestRegistry_ListSinks(t *testing.T) {
	r := adapter.NewRegistry()

	r.RegisterSink(&mockSinkAdapter{name: "sink-a"})
	r.RegisterSink(&mockSinkAdapter{name: "sink-b"})

	names := r.ListSinks()
	if len(names) != 2 {
		t.Errorf("expected 2 sinks, got %d", len(names))
	}
}
