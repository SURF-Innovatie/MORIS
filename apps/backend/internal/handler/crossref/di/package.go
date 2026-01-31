package di

import (
	"github.com/SURF-Innovatie/MORIS/internal/app/doi"
	doihandler "github.com/SURF-Innovatie/MORIS/internal/handler/doi"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideDoiHandler),
)

func provideDoiHandler(i do.Injector) (*doihandler.Handler, error) {
	svc := do.MustInvoke[doi.Service](i)
	return doihandler.NewHandler(svc), nil
}
