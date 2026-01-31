package di

import (
	"github.com/SURF-Innovatie/MORIS/internal/app/portfolio"
	portfoliorepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/portfolio"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(providePortfolioService),
)

func providePortfolioService(i do.Injector) (portfolio.Service, error) {
	repo := do.MustInvoke[*portfoliorepo.EntRepo](i)
	return portfolio.NewService(repo), nil
}
