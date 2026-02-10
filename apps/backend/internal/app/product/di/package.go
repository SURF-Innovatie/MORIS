package di

import (
	"github.com/SURF-Innovatie/MORIS/internal/app/doi"
	"github.com/SURF-Innovatie/MORIS/internal/app/product"
	productrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/product"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideProductService),
)

func provideProductService(i do.Injector) (product.Service, error) {
	doiSvc := do.MustInvoke[doi.Service](i)
	repo := do.MustInvoke[*productrepo.EntRepo](i)
	return product.NewService(repo, doiSvc), nil
}
