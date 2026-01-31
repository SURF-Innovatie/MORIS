package di

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	portfoliorepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/portfolio"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(providePortfolioRepo),
)

func providePortfolioRepo(i do.Injector) (*portfoliorepo.EntRepo, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return portfoliorepo.NewEntRepo(cli), nil
}
