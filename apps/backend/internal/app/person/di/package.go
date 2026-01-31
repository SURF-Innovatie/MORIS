package di

import (
	"github.com/SURF-Innovatie/MORIS/internal/app/person"
	personrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/person"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(ProvideService),
)

func ProvideService(i do.Injector) (person.Service, error) {
	repo := do.MustInvoke[*personrepo.EntRepo](i)
	return person.NewService(repo), nil
}
