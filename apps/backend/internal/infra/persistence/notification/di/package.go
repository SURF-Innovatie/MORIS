package di

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	notificationrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/notification"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideNotificationRepo),
)

func provideNotificationRepo(i do.Injector) (*notificationrepo.EntRepo, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return notificationrepo.NewEntRepo(cli), nil
}
