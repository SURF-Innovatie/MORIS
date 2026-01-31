package di

import (
	"github.com/SURF-Innovatie/MORIS/internal/app/doi"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideDoiService),
)

func provideDoiService(i do.Injector) (doi.Service, error) {
	return doi.NewService(), nil
}
