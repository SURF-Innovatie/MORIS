package eventstore

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
)

// Stream stores events per aggregate for tests / in-memory runs.
type Stream map[uuid.UUID][]events.Event

func NewStream() Stream { return make(Stream) }

func (s Stream) append(id uuid.UUID, ev events.Event) {
	s[id] = append(s[id], ev)
}

func (s Stream) events(id uuid.UUID) []events.Event {
	return append([]events.Event(nil), s[id]...)
}

// MemoryStore implements Store using Stream.
type MemoryStore struct {
	mu     sync.RWMutex
	stream Stream
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		stream: NewStream(),
	}
}

func (m *MemoryStore) Load(ctx context.Context, id uuid.UUID) ([]events.Event, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.stream.events(id), nil
}

func (m *MemoryStore) Append(ctx context.Context, id uuid.UUID, newEvents ...events.Event) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, ev := range newEvents {
		m.stream.append(id, ev)
	}
	return nil
}
