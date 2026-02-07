package di

import (
	appnotification "github.com/SURF-Innovatie/MORIS/internal/app/notification"
	"github.com/SURF-Innovatie/MORIS/internal/handler/ldn/inbox"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideLDNInboxHandler),
)

func provideLDNInboxHandler(i do.Injector) (*inbox.Handler, error) {
	svc := do.MustInvoke[appnotification.Service](i)
	return inbox.NewHandler(svc), nil
}
