package di

import (
	"github.com/SURF-Innovatie/MORIS/internal/app/eventpolicy"
	eventpolicyhandler "github.com/SURF-Innovatie/MORIS/internal/handler/eventpolicy"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideEventPolicyHandler),
)

func provideEventPolicyHandler(i do.Injector) (*eventpolicyhandler.Handler, error) {
	svc := do.MustInvoke[eventpolicy.Service](i)
	return eventpolicyhandler.NewHandler(svc), nil
}
