package di

import (
	authapp "github.com/SURF-Innovatie/MORIS/internal/app/auth"
	"github.com/SURF-Innovatie/MORIS/internal/infra/env"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(ProvideService),
)

func ProvideService(i do.Injector) (authapp.Service, error) {
	repo := do.MustInvoke[authapp.Repository](i)
	return authapp.NewService(repo, authapp.Options{
		JWTSecret: env.Global.JWTSecret,
	}), nil
}
