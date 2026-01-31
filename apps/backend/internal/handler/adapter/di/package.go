package di

import (
	"github.com/SURF-Innovatie/MORIS/internal/adapter"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/queries"
	adapterhandler "github.com/SURF-Innovatie/MORIS/internal/handler/adapter"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideAdapterHandler),
)

func provideAdapterHandler(i do.Injector) (*adapterhandler.Handler, error) {
	registry := do.MustInvoke[*adapter.Registry](i)
	projSvc := do.MustInvoke[queries.Service](i)
	return adapterhandler.NewHandler(registry, projSvc), nil
}
