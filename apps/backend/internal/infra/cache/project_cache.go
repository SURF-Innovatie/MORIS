package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	ErrCacheNotInitialized = errors.New("redis not initialized")
	ErrCacheMiss           = errors.New("cache miss")
)

type ProjectCache interface {
	GetProject(ctx context.Context, id uuid.UUID) (*entities.Project, error)
	SetProject(ctx context.Context, proj *entities.Project) error
	DeleteProject(ctx context.Context, id uuid.UUID) error
}

type RedisProjectCache struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewRedisProjectCache(rdb *redis.Client, ttl time.Duration) *RedisProjectCache {
	return &RedisProjectCache{rdb: rdb, ttl: ttl}
}

func (c *RedisProjectCache) key(id uuid.UUID) string {
	return fmt.Sprintf("project:%s", id.String())
}

func (c *RedisProjectCache) GetProject(ctx context.Context, id uuid.UUID) (*entities.Project, error) {
	if c.rdb == nil {
		return nil, ErrCacheNotInitialized
	}

	val, err := c.rdb.Get(ctx, c.key(id)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrCacheMiss
		}
		return nil, err
	}

	var proj entities.Project
	if err := json.Unmarshal([]byte(val), &proj); err != nil {
		return nil, err
	}

	return &proj, nil
}

func (c *RedisProjectCache) SetProject(ctx context.Context, proj *entities.Project) error {
	if c.rdb == nil {
		return ErrCacheNotInitialized
	}
	if proj == nil {
		return errors.New("cannot cache nil project")
	}

	b, err := json.Marshal(proj)
	if err != nil {
		return err
	}

	return c.rdb.Set(ctx, c.key(proj.Id), b, c.ttl).Err()
}

func (c *RedisProjectCache) DeleteProject(ctx context.Context, id uuid.UUID) error {
	if c.rdb == nil {
		return nil
	}
	return c.rdb.Del(ctx, c.key(id)).Err()
}
