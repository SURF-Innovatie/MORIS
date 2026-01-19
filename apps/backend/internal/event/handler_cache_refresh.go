package event

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/SURF-Innovatie/MORIS/internal/infra/cache"
)

// CacheRefreshHandler refreshes project cache when events status changes
type CacheRefreshHandler struct {
	refresher cache.ProjectCacheRefresher
}

func NewCacheRefreshHandler(refresher cache.ProjectCacheRefresher) *CacheRefreshHandler {
	return &CacheRefreshHandler{refresher: refresher}
}

func (h *CacheRefreshHandler) Handle(ctx context.Context, e events.Event) error {
	_, err := h.refresher.Refresh(ctx, e.AggregateID())
	return err
}
