package adapter

import (
	"fmt"
	"sync"
)

// Registry manages available source and sink adapters.
type Registry struct {
	mu      sync.RWMutex
	sources map[string]SourceAdapter
	sinks   map[string]SinkAdapter
}

// NewRegistry creates a new adapter registry.
func NewRegistry() *Registry {
	return &Registry{
		sources: make(map[string]SourceAdapter),
		sinks:   make(map[string]SinkAdapter),
	}
}

// RegisterSource registers a source adapter.
func (r *Registry) RegisterSource(adapter SourceAdapter) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := adapter.Name()
	if _, exists := r.sources[name]; exists {
		return fmt.Errorf("source adapter %q already registered", name)
	}
	r.sources[name] = adapter
	return nil
}

// RegisterSink registers a sink adapter.
func (r *Registry) RegisterSink(adapter SinkAdapter) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := adapter.Name()
	if _, exists := r.sinks[name]; exists {
		return fmt.Errorf("sink adapter %q already registered", name)
	}
	r.sinks[name] = adapter
	return nil
}

// GetSource retrieves a source adapter by name.
func (r *Registry) GetSource(name string) (SourceAdapter, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	adapter, ok := r.sources[name]
	return adapter, ok
}

// GetSink retrieves a sink adapter by name.
func (r *Registry) GetSink(name string) (SinkAdapter, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	adapter, ok := r.sinks[name]
	return adapter, ok
}

// ListSources returns the names of all registered source adapters.
func (r *Registry) ListSources() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.sources))
	for name := range r.sources {
		names = append(names, name)
	}
	return names
}

// ListSinks returns the names of all registered sink adapters.
func (r *Registry) ListSinks() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.sinks))
	for name := range r.sinks {
		names = append(names, name)
	}
	return names
}
