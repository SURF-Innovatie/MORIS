package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrCacheNotInitialized = errors.New("redis not initialized")
	ErrCacheMiss           = errors.New("cache miss")
)

type Cache[T any] interface {
	Get(ctx context.Context, key string) (*T, error)
	Set(ctx context.Context, key string, value *T) error
	Delete(ctx context.Context, key string) error
}

type RedisCache[T any] struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewRedisCache[T any](rdb *redis.Client, ttl time.Duration) *RedisCache[T] {
	return &RedisCache[T]{rdb: rdb, ttl: ttl}
}

func (c *RedisCache[T]) Get(ctx context.Context, key string) (*T, error) {
	if c.rdb == nil {
		return nil, ErrCacheNotInitialized
	}

	val, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrCacheMiss
		}
		return nil, err
	}

	var item T
	if err := json.Unmarshal([]byte(val), &item); err != nil {
		return nil, err
	}

	return &item, nil
}

func (c *RedisCache[T]) Set(ctx context.Context, key string, value *T) error {
	if c.rdb == nil {
		return ErrCacheNotInitialized
	}
	if value == nil {
		return errors.New("cannot cache nil value")
	}

	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.rdb.Set(ctx, key, b, c.ttl).Err()
}

func (c *RedisCache[T]) Delete(ctx context.Context, key string) error {
	if c.rdb == nil {
		return nil
	}
	return c.rdb.Del(ctx, key).Err()
}
