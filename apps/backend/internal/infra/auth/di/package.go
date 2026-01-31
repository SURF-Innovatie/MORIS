package di

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	coreauth "github.com/SURF-Innovatie/MORIS/internal/app/auth"
	"github.com/SURF-Innovatie/MORIS/internal/infra/auth"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideCurrentUserProvider),
)

func provideCurrentUserProvider(i do.Injector) (coreauth.CurrentUserProvider, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return auth.NewCurrentUserProvider(cli), nil
}
