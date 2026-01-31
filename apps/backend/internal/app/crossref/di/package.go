package di

import (
	excrossref "github.com/SURF-Innovatie/MORIS/external/crossref"
	"github.com/SURF-Innovatie/MORIS/internal/app/crossref"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideCrossrefService),
)

func provideCrossrefService(i do.Injector) (crossref.Service, error) {
	cli := do.MustInvoke[excrossref.Client](i)
	return crossref.NewService(cli), nil
}
