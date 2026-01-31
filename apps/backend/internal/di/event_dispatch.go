package di

import (
	"github.com/SURF-Innovatie/MORIS/internal/app/event"
	"github.com/SURF-Innovatie/MORIS/internal/infra/eventdispatch"
	"github.com/SURF-Innovatie/MORIS/internal/infra/handlers/events"
	"github.com/samber/do/v2"
)

func provideEventPublisher(i do.Injector) (event.Publisher, error) {
	projHandler := do.MustInvoke[*events.ProjectEventNotificationHandler](i)
	approvalHandler := do.MustInvoke[*events.ApprovalRequestNotificationHandler](i)
	policyHandler := do.MustInvoke[*events.Handler](i)
	execHandler := do.MustInvoke[*events.PolicyExecutionHandler](i)

	notificationHandlers := []eventdispatch.NotificationHandler{
		projHandler,
		approvalHandler,
		policyHandler,
		execHandler,
	}

	statusHandler := do.MustInvoke[*events.StatusUpdateNotificationHandler](i)
	cacheHandler := do.MustInvoke[*events.CacheRefreshHandler](i)

	statusChangeHandlers := []eventdispatch.StatusChangeHandler{
		statusHandler.Handle,
		cacheHandler.Handle,
	}

	return eventdispatch.New(notificationHandlers, statusChangeHandlers), nil
}
