package di

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/app/catalog"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/queries"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideService),
	do.Lazy(provideHandler),
)

func provideService(i do.Injector) (catalog.Service, error) {
	client := do.MustInvoke[*ent.Client](i)
	projectSvc := do.MustInvoke[queries.Service](i)
	return catalog.NewService(client, projectSvc), nil
}

func provideHandler(i do.Injector) (*catalog.Handler, error) {
	svc := do.MustInvoke[catalog.Service](i)
	return catalog.NewHandler(svc), nil
}
