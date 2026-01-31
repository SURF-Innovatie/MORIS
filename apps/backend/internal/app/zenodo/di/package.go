package di

import (
	exzenodo "github.com/SURF-Innovatie/MORIS/external/zenodo"
	"github.com/SURF-Innovatie/MORIS/internal/app/zenodo"
	userrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/user"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideZenodoService),
)

func provideZenodoService(i do.Injector) (zenodo.Service, error) {
	userRepo := do.MustInvoke[*userrepo.EntRepo](i)
	cli := do.MustInvoke[*exzenodo.Client](i)
	return zenodo.NewService(userRepo, cli), nil
}
