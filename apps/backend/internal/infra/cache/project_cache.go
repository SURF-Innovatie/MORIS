package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type ProjectCache interface {
	GetProject(ctx context.Context, id uuid.UUID) (*entities.Project, error)
	SetProject(ctx context.Context, proj *entities.Project) error
	DeleteProject(ctx context.Context, id uuid.UUID) error
}

type RedisProjectCache struct {
	*RedisCache[entities.Project]
}

func NewRedisProjectCache(rdb *redis.Client, ttl time.Duration) *RedisProjectCache {
	return &RedisProjectCache{
		RedisCache: NewRedisCache[entities.Project](rdb, ttl),
	}
}

func (c *RedisProjectCache) key(id uuid.UUID) string {
	return fmt.Sprintf("project:%s", id.String())
}

func (c *RedisProjectCache) GetProject(ctx context.Context, id uuid.UUID) (*entities.Project, error) {
	return c.Get(ctx, c.key(id))
}

func (c *RedisProjectCache) SetProject(ctx context.Context, proj *entities.Project) error {
	if proj == nil {
		// Preserve original error message if preferred, or rely on generic
		// Generic says "cannot cache nil value"
		return fmt.Errorf("cannot cache nil project")
	}
	// Note: Generic Set checks for nil value too, but we check here for custom message or safety BEFORE key generation (though key generation here is safe if we used proj.Id which would panic if proj is nil).
	// Actually original code used `proj.Id` effectively.
	return c.Set(ctx, c.key(proj.Id), proj)
}

func (c *RedisProjectCache) DeleteProject(ctx context.Context, id uuid.UUID) error {
	return c.Delete(ctx, c.key(id))
}
