package di

import (
	coreauth "github.com/SURF-Innovatie/MORIS/internal/app/auth"
	"github.com/SURF-Innovatie/MORIS/internal/app/product"
	producthandler "github.com/SURF-Innovatie/MORIS/internal/handler/product"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideProductHandler),
)

func provideProductHandler(i do.Injector) (*producthandler.Handler, error) {
	svc := do.MustInvoke[product.Service](i)
	curUser := do.MustInvoke[coreauth.CurrentUserProvider](i)
	return producthandler.NewHandler(svc, curUser), nil
}
