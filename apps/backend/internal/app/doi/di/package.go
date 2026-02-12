package di

import (
	doiex "github.com/SURF-Innovatie/MORIS/external/doi"
	"github.com/SURF-Innovatie/MORIS/internal/app/doi"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideDoiService),
)

func provideDoiService(i do.Injector) (doi.Service, error) {
	doiClient := do.MustInvoke[doiex.Client](i)
	return doi.NewService(doiClient), nil
}
