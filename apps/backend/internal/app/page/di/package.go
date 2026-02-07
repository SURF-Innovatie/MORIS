package di

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/app/page"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(ProvideRepository),
	do.Lazy(ProvideService),
)

func ProvideRepository(i do.Injector) (page.Repository, error) {
	client := do.MustInvoke[*ent.Client](i)
	return page.NewRepository(client), nil
}

func ProvideService(i do.Injector) (page.Service, error) {
	repo := do.MustInvoke[page.Repository](i)
	return page.NewService(repo), nil
}
