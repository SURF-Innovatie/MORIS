package di

import (
	"github.com/SURF-Innovatie/MORIS/internal/app/event"
	"github.com/SURF-Innovatie/MORIS/internal/infra/eventdispatch"
	"github.com/SURF-Innovatie/MORIS/internal/infra/handlers/events"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideEventPublisher),
)

func provideEventPublisher(i do.Injector) (event.Publisher, error) {
	policyHandler := do.MustInvoke[*events.Handler](i)
	execHandler := do.MustInvoke[*events.PolicyExecutionHandler](i)

	notificationHandlers := []eventdispatch.NotificationHandler{
		policyHandler,
		execHandler,
	}

	cacheHandler := do.MustInvoke[*events.CacheRefreshHandler](i)

	statusChangeHandlers := []eventdispatch.StatusChangeHandler{
		cacheHandler.Handle,
	}

	return eventdispatch.New(notificationHandlers, statusChangeHandlers), nil
}
