package di

import (
	exnwo "github.com/SURF-Innovatie/MORIS/external/nwo"
	"github.com/SURF-Innovatie/MORIS/internal/app/nwo"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideNWOService),
)

func provideNWOService(i do.Injector) (nwo.Service, error) {
	cli := do.MustInvoke[exnwo.Client](i)
	return nwo.NewService(cli), nil
}
