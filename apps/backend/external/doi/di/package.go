package di

import (
	"github.com/SURF-Innovatie/MORIS/external/doi"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideDoiClient),
)

func provideDoiClient(i do.Injector) (doi.Client, error) {
	return doi.NewClient(), nil
}
