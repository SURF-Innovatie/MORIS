package di

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/app/eventpolicy"
	personsvc "github.com/SURF-Innovatie/MORIS/internal/app/person"
	personhandler "github.com/SURF-Innovatie/MORIS/internal/handler/person"
	"github.com/SURF-Innovatie/MORIS/internal/infra/cache"
	"github.com/SURF-Innovatie/MORIS/internal/infra/handlers/events"
	eventrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/event"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideEventPolicyHandler),
	do.Lazy(providePolicyExecutionHandler),
	do.Lazy(provideCacheRefreshHandler),
)

func providePolicyExecutionHandler(i do.Injector) (*events.PolicyExecutionHandler, error) {
	evaluator := do.MustInvoke[eventpolicy.Evaluator](i)
	entRepo := do.MustInvoke[*eventrepo.EntRepo](i)
	return events.NewPolicyExecutionHandler(evaluator, entRepo), nil
}

func provideEventPolicyHandler(i do.Injector) (*events.Handler, error) {
	repo := do.MustInvoke[eventpolicy.Service](i)
	cli := do.MustInvoke[*ent.Client](i)
	return events.NewEventPolicyHandler(repo, cli), nil
}

func provideCacheRefreshHandler(i do.Injector) (*events.CacheRefreshHandler, error) {
	refresher := do.MustInvoke[cache.ProjectCacheRefresher](i)
	return events.NewCacheRefreshHandler(refresher), nil
}

// HTTP Handler Providers

func providePersonHandler(i do.Injector) (*personhandler.Handler, error) {
	svc := do.MustInvoke[personsvc.Service](i)
	return personhandler.NewHandler(svc), nil
}
