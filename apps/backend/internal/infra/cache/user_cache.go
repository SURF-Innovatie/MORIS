package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type UserCache interface {
	GetUser(ctx context.Context, id uuid.UUID) (*entities.User, error)
	SetUser(ctx context.Context, user *entities.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type RedisUserCache struct {
	*RedisCache[entities.User]
}

func NewRedisUserCache(rdb *redis.Client, ttl time.Duration) *RedisUserCache {
	return &RedisUserCache{
		RedisCache: NewRedisCache[entities.User](rdb, ttl),
	}
}

func (c *RedisUserCache) key(id uuid.UUID) string {
	return fmt.Sprintf("user:%s", id.String())
}

func (c *RedisUserCache) GetUser(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	return c.Get(ctx, c.key(id))
}

func (c *RedisUserCache) SetUser(ctx context.Context, user *entities.User) error {
	if user == nil {
		return fmt.Errorf("cannot cache nil user")
	}
	return c.Set(ctx, c.key(user.ID), user)
}

func (c *RedisUserCache) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return c.Delete(ctx, c.key(id))
}
