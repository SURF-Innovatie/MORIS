package di

import (
	personsvc "github.com/SURF-Innovatie/MORIS/internal/app/person"
	personhandler "github.com/SURF-Innovatie/MORIS/internal/handler/person"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(ProvideHandler),
)

func ProvideHandler(i do.Injector) (*personhandler.Handler, error) {
	svc := do.MustInvoke[personsvc.Service](i)
	return personhandler.NewHandler(svc), nil
}
