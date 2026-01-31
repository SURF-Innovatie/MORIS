package di

import (
	exsurfconext "github.com/SURF-Innovatie/MORIS/external/surfconext"
	coreauth "github.com/SURF-Innovatie/MORIS/internal/app/auth"
	"github.com/SURF-Innovatie/MORIS/internal/app/surfconext"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideSurfconextService),
)

func provideSurfconextService(i do.Injector) (surfconext.Service, error) {
	cli := do.MustInvoke[*exsurfconext.Client](i)
	authSvc := do.MustInvoke[coreauth.Service](i)
	return surfconext.NewService(cli, authSvc), nil
}
