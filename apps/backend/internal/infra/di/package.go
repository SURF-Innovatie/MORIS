package di

import (
	"context"
	"fmt"
	"time"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/command"
	"github.com/SURF-Innovatie/MORIS/internal/infra/cache"
	"github.com/SURF-Innovatie/MORIS/internal/infra/env"
	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/entclient"
	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/enttx"
	eventrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/event"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideEntClient),
	do.Lazy(provideRedisClient),

	do.Lazy(provideProjectCache),
	do.Lazy(provideUserCache),
	do.Lazy(provideCacheRefresher),

	// tx, ent providers, etc
	do.Lazy(provideTxManager),
	do.Lazy(provideEntClientProvider),
)

func provideEntClient(i do.Injector) (*ent.Client, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		env.Global.DBHost, env.Global.DBPort, env.Global.DBUser, env.Global.DBPassword, env.Global.DBName)
	return ent.Open("postgres", dsn)
}

func provideRedisClient(i do.Injector) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", env.Global.CacheHost, env.Global.CachePort),
		Password: env.Global.CachePassword,
		Username: env.Global.CacheUser,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Warn().Err(err).Msg("failed to connect to redis/valkey")
	} else {
		log.Info().Msgf("Connected to Redis at %s:%s", env.Global.CacheHost, env.Global.CachePort)
	}
	return rdb, nil
}

func provideProjectCache(i do.Injector) (cache.ProjectCache, error) {
	rdb := do.MustInvoke[*redis.Client](i)
	return cache.NewRedisProjectCache(rdb, 24*time.Hour), nil
}

func provideUserCache(i do.Injector) (cache.UserCache, error) {
	rdb := do.MustInvoke[*redis.Client](i)
	return cache.NewRedisUserCache(rdb, 1*time.Hour), nil
}

func provideCacheRefresher(i do.Injector) (cache.ProjectCacheRefresher, error) {
	eventRepo := do.MustInvoke[*eventrepo.EntRepo](i)
	pc := do.MustInvoke[cache.ProjectCache](i)
	return cache.NewEventStoreProjectCacheRefresher(eventRepo, pc), nil
}

func provideEntClientProvider(i do.Injector) (command.EntClientProvider, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return entclient.New(cli), nil
}

func provideTxManager(i do.Injector) (*enttx.Manager, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return enttx.NewManager(cli), nil
}
