package di

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	coreauth "github.com/SURF-Innovatie/MORIS/internal/app/auth"
	personsvc "github.com/SURF-Innovatie/MORIS/internal/app/person"
	"github.com/SURF-Innovatie/MORIS/internal/app/user"
	"github.com/SURF-Innovatie/MORIS/internal/infra/auth"
	"github.com/SURF-Innovatie/MORIS/internal/infra/env"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideAuthService),
)

func provideAuthService(i do.Injector) (coreauth.Service, error) {
	cli := do.MustInvoke[*ent.Client](i)
	userSvc := do.MustInvoke[user.Service](i)
	personSvc := do.MustInvoke[personsvc.Service](i)
	return auth.NewJWTService(cli, userSvc, personSvc, env.Global.JWTSecret), nil
}
