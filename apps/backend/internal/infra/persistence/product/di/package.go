package di

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	productrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/product"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideProductRepo),
)

func provideProductRepo(i do.Injector) (*productrepo.EntRepo, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return productrepo.NewEntRepo(cli), nil
}
