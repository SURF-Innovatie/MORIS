package di

import (
	"github.com/SURF-Innovatie/MORIS/internal/app/nwo"
	nwohandler "github.com/SURF-Innovatie/MORIS/internal/handler/nwo"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideNWOHandler),
)

func provideNWOHandler(i do.Injector) (*nwohandler.Handler, error) {
	svc := do.MustInvoke[nwo.Service](i)
	return nwohandler.NewHandler(svc), nil
}
