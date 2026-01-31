package di

import (
	"github.com/SURF-Innovatie/MORIS/internal/app/crossref"
	crossrefhandler "github.com/SURF-Innovatie/MORIS/internal/handler/crossref"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideCrossrefHandler),
)

func provideCrossrefHandler(i do.Injector) (*crossrefhandler.Handler, error) {
	svc := do.MustInvoke[crossref.Service](i)
	return crossrefhandler.NewHandler(svc), nil
}
