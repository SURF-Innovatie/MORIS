package cachewarmup

import (
	"context"
	"errors"

	"github.com/SURF-Innovatie/MORIS/internal/app/project/load"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/queries"
	"github.com/SURF-Innovatie/MORIS/internal/infra/cache"
)

type Service interface {
	WarmupProjects(ctx context.Context) (int, error) // returns how many cached
}

type service struct {
	repo   queries.ProjectReadRepository
	loader *load.Loader
	cache  cache.ProjectCache
}

func NewService(repo queries.ProjectReadRepository, loader *load.Loader, c cache.ProjectCache) Service {
	return &service{
		repo:   repo,
		loader: loader,
		cache:  c,
	}
}

func (s *service) WarmupProjects(ctx context.Context) (int, error) {
	if s.cache == nil {
		return 0, cache.ErrCacheNotInitialized
	}
	if s.loader == nil {
		return 0, errors.New("nil loader")
	}

	ids, err := s.repo.ProjectIDsStarted(ctx)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, id := range ids {
		proj, err := s.loader.Load(ctx, id)
		if err != nil {
			// skip broken projects during warmup
			continue
		}
		if err := s.cache.SetProject(ctx, proj); err != nil {
			// skip cache write errors, keep going
			continue
		}
		count++
	}

	return count, nil
}
