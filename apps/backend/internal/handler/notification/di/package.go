package di

import (
	"github.com/SURF-Innovatie/MORIS/internal/app/notification"
	notificationhandler "github.com/SURF-Innovatie/MORIS/internal/handler/notification"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideNotificationHandler),
)

func provideNotificationHandler(i do.Injector) (*notificationhandler.Handler, error) {
	svc := do.MustInvoke[notification.Service](i)
	return notificationhandler.NewHandler(svc), nil
}
