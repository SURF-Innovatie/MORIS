package di

import (
	"github.com/SURF-Innovatie/MORIS/internal/app/page"
	pagehandler "github.com/SURF-Innovatie/MORIS/internal/handler/page"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(ProvideHandler),
)

func ProvideHandler(i do.Injector) (*pagehandler.Handler, error) {
	service := do.MustInvoke[page.Service](i)
	return pagehandler.NewHandler(service), nil
}
