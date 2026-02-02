package di

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	authapp "github.com/SURF-Innovatie/MORIS/internal/app/auth"
	"github.com/SURF-Innovatie/MORIS/internal/infra/identity"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(ProvideCurrentUserProvider),
)

func ProvideCurrentUserProvider(i do.Injector) (authapp.CurrentUserProvider, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return identity.NewCurrentUserProvider(cli), nil
}
