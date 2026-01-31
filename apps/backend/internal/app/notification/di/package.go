package di

import (
	"github.com/SURF-Innovatie/MORIS/internal/app/notification"
	notificationrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/notification"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideNotificationService),
)

func provideNotificationService(i do.Injector) (notification.Service, error) {
	repo := do.MustInvoke[*notificationrepo.EntRepo](i)
	return notification.NewService(repo), nil
}
