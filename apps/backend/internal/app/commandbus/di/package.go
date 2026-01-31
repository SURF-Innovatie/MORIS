package di

import (
	"github.com/SURF-Innovatie/MORIS/internal/app/event"
	"github.com/SURF-Innovatie/MORIS/internal/app/notification"
	eventrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/event"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideEventService),
)

func provideEventService(i do.Injector) (event.Service, error) {
	notifSvc := do.MustInvoke[notification.Service](i)
	repo := do.MustInvoke[*eventrepo.EntRepo](i)

	evtPub := do.MustInvoke[event.Publisher](i)

	return event.NewService(repo, notifSvc, evtPub), nil
}
