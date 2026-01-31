package di

import (
	coreauth "github.com/SURF-Innovatie/MORIS/internal/app/auth"
	"github.com/SURF-Innovatie/MORIS/internal/app/zenodo"
	zenodohandler "github.com/SURF-Innovatie/MORIS/internal/handler/zenodo"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideZenodoHandler),
)

func provideZenodoHandler(i do.Injector) (*zenodohandler.Handler, error) {
	svc := do.MustInvoke[zenodo.Service](i)
	curUser := do.MustInvoke[coreauth.CurrentUserProvider](i)
	return zenodohandler.NewHandler(svc, curUser), nil
}
