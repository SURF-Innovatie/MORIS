package di

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	userrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/user"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideUserRepo),
)

func provideUserRepo(i do.Injector) (*userrepo.EntRepo, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return userrepo.NewEntRepo(cli), nil
}
