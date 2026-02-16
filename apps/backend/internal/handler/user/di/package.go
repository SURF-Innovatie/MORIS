package di

import (
	"github.com/SURF-Innovatie/MORIS/internal/app/event"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/queries"
	"github.com/SURF-Innovatie/MORIS/internal/app/user"
	"github.com/SURF-Innovatie/MORIS/internal/domain/project/events/hydrator"
	userhandler "github.com/SURF-Innovatie/MORIS/internal/handler/user"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideUserHandler),
)

func provideUserHandler(i do.Injector) (*userhandler.Handler, error) {
	userSvc := do.MustInvoke[user.Service](i)
	projSvc := do.MustInvoke[queries.Service](i)
	eventSvc := do.MustInvoke[event.Service](i)
	h := do.MustInvoke[*hydrator.Hydrator](i)
	return userhandler.NewHandler(userSvc, projSvc, eventSvc, h), nil
}
