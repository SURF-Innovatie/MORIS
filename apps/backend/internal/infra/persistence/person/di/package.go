package di

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	personrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/person"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(ProvideRepo),
)

func ProvideRepo(i do.Injector) (*personrepo.EntRepo, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return personrepo.NewEntRepo(cli), nil
}
