package cache

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/app/commandbus"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/projection"
	"github.com/google/uuid"
)

type ProjectCacheRefresher interface {
	Refresh(ctx context.Context, projectID uuid.UUID) (*entities.Project, error)
}

type EventStoreProjectCacheRefresher struct {
	eventStore commandbus.EventStore
	cache      ProjectCache
}

func NewEventStoreProjectCacheRefresher(eventStore commandbus.EventStore, cache ProjectCache) *EventStoreProjectCacheRefresher {
	return &EventStoreProjectCacheRefresher{eventStore: eventStore, cache: cache}
}

func (r *EventStoreProjectCacheRefresher) Refresh(ctx context.Context, projectID uuid.UUID) (*entities.Project, error) {
	evts, version, err := r.eventStore.Load(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if len(evts) == 0 {
		return nil, nil
	}

	proj := projection.Reduce(projectID, evts)
	if proj == nil {
		return nil, nil
	}
	proj.Version = version

	_ = r.cache.SetProject(ctx, proj)
	return proj, nil
}
