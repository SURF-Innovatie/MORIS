package load

import (
	"context"
	"errors"

	"github.com/SURF-Innovatie/MORIS/internal/app/event"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/projection"
	"github.com/google/uuid"
)

var ErrNotFound = errors.New("project not found")

type Loader struct {
	eventSvc event.Service
	cache    Cache
}

func New(eventSvc event.Service, cache Cache) *Loader {
	return &Loader{eventSvc: eventSvc, cache: cache}
}

func (l *Loader) Load(ctx context.Context, projectID uuid.UUID) (*entities.Project, error) {
	if l.cache != nil {
		if p, err := l.cache.GetProject(ctx, projectID); err == nil && p != nil {
			return p, nil
		}
	}

	evts, version, err := l.eventSvc.Load(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if len(evts) == 0 {
		return nil, ErrNotFound
	}

	p := projection.Reduce(projectID, evts)
	p.Version = version

	_ = l.cache.SetProject(ctx, p)
	return p, nil
}
