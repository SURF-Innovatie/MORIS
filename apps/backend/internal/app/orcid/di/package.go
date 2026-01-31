package di

import (
	exorcid "github.com/SURF-Innovatie/MORIS/external/orcid"
	"github.com/SURF-Innovatie/MORIS/internal/app/orcid"
	personrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/person"
	userrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/user"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideORCIDService),
)

func provideORCIDService(i do.Injector) (orcid.Service, error) {
	userRepo := do.MustInvoke[*userrepo.EntRepo](i)
	personRepo := do.MustInvoke[*personrepo.EntRepo](i)
	cli := do.MustInvoke[*exorcid.Client](i)
	return orcid.NewService(userRepo, personRepo, cli), nil
}
