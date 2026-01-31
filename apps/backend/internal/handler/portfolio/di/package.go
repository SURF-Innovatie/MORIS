package di

import (
	coreauth "github.com/SURF-Innovatie/MORIS/internal/app/auth"
	"github.com/SURF-Innovatie/MORIS/internal/app/portfolio"
	portfoliohandler "github.com/SURF-Innovatie/MORIS/internal/handler/portfolio"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(providePortfolioHandler),
)

func providePortfolioHandler(i do.Injector) (*portfoliohandler.Handler, error) {
	svc := do.MustInvoke[portfolio.Service](i)
	curUser := do.MustInvoke[coreauth.CurrentUserProvider](i)
	return portfoliohandler.NewHandler(svc, curUser), nil
}
