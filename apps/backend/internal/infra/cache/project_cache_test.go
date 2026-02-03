package cache_test

import (
	"context"
	"testing"
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/domain/project"
	"github.com/SURF-Innovatie/MORIS/internal/infra/cache"
	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func TestRedisProjectCache_SetAndGet(t *testing.T) {
	ctx := context.Background()

	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { _ = rdb.Close() })

	c := cache.NewRedisProjectCache(rdb, 24*time.Hour)

	id := uuid.New()
	p := &project.Project{
		Id:    id,
		Title: "My Project",
	}

	if err := c.SetProject(ctx, p); err != nil {
		t.Fatalf("SetProject() err = %v", err)
	}

	got, err := c.GetProject(ctx, id)
	if err != nil {
		t.Fatalf("GetProject() err = %v", err)
	}
	if got == nil {
		t.Fatalf("GetProject() got nil")
	}
	if got.Id != id {
		t.Fatalf("GetProject() id mismatch: got %s want %s", got.Id, id)
	}
	if got.Title != "My Project" {
		t.Fatalf("GetProject() title mismatch: got %q want %q", got.Title, "My Project")
	}
}

func TestRedisProjectCache_GetMiss(t *testing.T) {
	ctx := context.Background()

	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { _ = rdb.Close() })

	c := cache.NewRedisProjectCache(rdb, 24*time.Hour)

	_, err := c.GetProject(ctx, uuid.New())
	if err == nil {
		t.Fatalf("GetProject() expected error, got nil")
	}
	if err != cache.ErrCacheMiss {
		t.Fatalf("GetProject() expected ErrCacheMiss, got %v", err)
	}
}

func TestRedisProjectCache_Delete(t *testing.T) {
	ctx := context.Background()

	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { _ = rdb.Close() })

	c := cache.NewRedisProjectCache(rdb, 24*time.Hour)

	id := uuid.New()
	p := &project.Project{Id: id, Title: "To be deleted"}

	if err := c.SetProject(ctx, p); err != nil {
		t.Fatalf("SetProject() err = %v", err)
	}

	if err := c.DeleteProject(ctx, id); err != nil {
		t.Fatalf("DeleteProject() err = %v", err)
	}

	_, err := c.GetProject(ctx, id)
	if err != cache.ErrCacheMiss {
		t.Fatalf("GetProject() expected ErrCacheMiss after delete, got %v", err)
	}
}

func TestRedisProjectCache_NilRedis(t *testing.T) {
	ctx := context.Background()

	pc := cache.NewRedisProjectCache(nil, time.Hour)

	projectID := uuid.New()
	proj := &project.Project{Id: projectID, Title: "x"}

	err := pc.SetProject(ctx, proj)
	if err == nil {
		t.Fatalf("SetProject() expected error, got nil")
	}
}
