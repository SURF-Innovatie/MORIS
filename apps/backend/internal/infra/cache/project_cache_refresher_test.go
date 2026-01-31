package cache_test

import (
	"context"
	"testing"
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/app/commandbus"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/SURF-Innovatie/MORIS/internal/infra/cache"
	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type fakeEventStore struct {
	evts    []events.Event
	version int
	err     error
}

func (f fakeEventStore) Load(ctx context.Context, id uuid.UUID) ([]events.Event, int, error) {
	return f.evts, f.version, f.err
}

func (f fakeEventStore) Append(ctx context.Context, id uuid.UUID, expectedVersion int, evts ...events.Event) error {
	return nil
}

func TestEventStoreProjectCacheRefresher_Refresh_WritesToCache(t *testing.T) {
	ctx := context.Background()

	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { _ = rdb.Close() })

	pc := cache.NewRedisProjectCache(rdb, 24*time.Hour)

	projectID := uuid.New()
	actorID := uuid.New()

	start := &events.ProjectStarted{
		Base:        events.NewBase(projectID, actorID, events.StatusApproved),
		Title:       "Test Title",
		Description: "Test Description",
	}

	es := fakeEventStore{
		evts:    []events.Event{start},
		version: 7,
		err:     nil,
	}

	// Compile-time assertion: fake implements the port we actually need
	var _ commandbus.EventStore = es

	ref := cache.NewEventStoreProjectCacheRefresher(es, pc)

	proj, err := ref.Refresh(ctx, projectID)
	if err != nil {
		t.Fatalf("Refresh() err = %v", err)
	}
	if proj == nil {
		t.Fatalf("Refresh() returned nil project")
	}
	if proj.Id != projectID {
		t.Fatalf("Refresh() id mismatch: got %s want %s", proj.Id, projectID)
	}
	if proj.Version != 7 {
		t.Fatalf("Refresh() version mismatch: got %d want %d", proj.Version, 7)
	}
	if proj.Title != "Test Title" {
		t.Fatalf("Refresh() title mismatch: got %q want %q", proj.Title, "Test Title")
	}

	got, err := pc.GetProject(ctx, projectID)
	if err != nil {
		t.Fatalf("cache.GetProject() err = %v", err)
	}
	if got.Version != 7 {
		t.Fatalf("cached project version mismatch: got %d want %d", got.Version, 7)
	}
	if got.Title != "Test Title" {
		t.Fatalf("cached project title mismatch: got %q want %q", got.Title, "Test Title")
	}
}

func TestEventStoreProjectCacheRefresher_Refresh_NoEvents_ReturnsNil(t *testing.T) {
	ctx := context.Background()

	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { _ = rdb.Close() })

	pc := cache.NewRedisProjectCache(rdb, 24*time.Hour)

	projectID := uuid.New()
	es := fakeEventStore{evts: nil, version: 0, err: nil}

	ref := cache.NewEventStoreProjectCacheRefresher(es, pc)

	proj, err := ref.Refresh(ctx, projectID)
	if err != nil {
		t.Fatalf("Refresh() err = %v", err)
	}
	if proj != nil {
		t.Fatalf("Refresh() expected nil project, got %v", proj)
	}

	_, err = pc.GetProject(ctx, projectID)
	if err != cache.ErrCacheMiss {
		t.Fatalf("expected cache miss, got %v", err)
	}
}
