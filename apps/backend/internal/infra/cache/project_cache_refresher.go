package cache

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/projection"
	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/eventstore"
	"github.com/google/uuid"
)

type ProjectCacheRefresher interface {
	Refresh(ctx context.Context, projectID uuid.UUID) (*entities.Project, error)
}

type EventstoreProjectCacheRefresher struct {
	es    eventstore.Store
	cache ProjectCache
}

func NewEventstoreProjectCacheRefresher(es eventstore.Store, cache ProjectCache) *EventstoreProjectCacheRefresher {
	return &EventstoreProjectCacheRefresher{es: es, cache: cache}
}

func (r *EventstoreProjectCacheRefresher) Refresh(ctx context.Context, projectID uuid.UUID) (*entities.Project, error) {
	evts, version, err := r.es.Load(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if len(evts) == 0 {
		return nil, nil
	}

	proj := projection.Reduce(projectID, evts)
	proj.Version = version

	_ = r.cache.SetProject(ctx, proj)
	return proj, nil
}
