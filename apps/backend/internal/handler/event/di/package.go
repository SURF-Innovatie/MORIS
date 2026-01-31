package di

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	event2 "github.com/SURF-Innovatie/MORIS/internal/app/event"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/queries"
	"github.com/SURF-Innovatie/MORIS/internal/app/user"
	eventhandler "github.com/SURF-Innovatie/MORIS/internal/handler/event"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideEventHandler),
)

func provideEventHandler(i do.Injector) (*eventhandler.Handler, error) {
	evtSvc := do.MustInvoke[event2.Service](i)
	projSvc := do.MustInvoke[queries.Service](i)
	userSvc := do.MustInvoke[user.Service](i)
	cli := do.MustInvoke[*ent.Client](i)
	return eventhandler.NewHandler(evtSvc, projSvc, userSvc, cli), nil
}
