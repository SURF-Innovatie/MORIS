package di

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	authapp "github.com/SURF-Innovatie/MORIS/internal/app/auth"
	"github.com/SURF-Innovatie/MORIS/internal/app/user"
	authrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/auth"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideAuthRepo),
)

func provideAuthRepo(i do.Injector) (authapp.Repository, error) {
	cli := do.MustInvoke[*ent.Client](i)
	userSvc := do.MustInvoke[user.Service](i)
	return authrepo.NewEntRepo(cli, userSvc), nil
}
