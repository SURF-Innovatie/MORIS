package di

import (
	"github.com/SURF-Innovatie/MORIS/internal/app/orcid"
	orcidhandler "github.com/SURF-Innovatie/MORIS/internal/handler/orcid"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideORCIDHandler),
)

func provideORCIDHandler(i do.Injector) (*orcidhandler.Handler, error) {
	svc := do.MustInvoke[orcid.Service](i)
	return orcidhandler.NewHandler(svc), nil
}
