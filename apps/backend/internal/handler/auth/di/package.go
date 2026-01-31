package di

import (
	coreauth "github.com/SURF-Innovatie/MORIS/internal/app/auth"
	"github.com/SURF-Innovatie/MORIS/internal/app/orcid"
	"github.com/SURF-Innovatie/MORIS/internal/app/surfconext"
	"github.com/SURF-Innovatie/MORIS/internal/app/user"
	authhandler "github.com/SURF-Innovatie/MORIS/internal/handler/auth"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideAuthHandler),
)

func provideAuthHandler(i do.Injector) (*authhandler.Handler, error) {
	userSvc := do.MustInvoke[user.Service](i)
	authSvc := do.MustInvoke[coreauth.Service](i)
	orcidSvc := do.MustInvoke[orcid.Service](i)
	surfSvc := do.MustInvoke[surfconext.Service](i)
	return authhandler.NewHandler(userSvc, authSvc, orcidSvc, surfSvc), nil
}
