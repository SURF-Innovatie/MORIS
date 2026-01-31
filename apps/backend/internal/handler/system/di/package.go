package di

import (
	systemhandler "github.com/SURF-Innovatie/MORIS/internal/handler/system"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideSystemHandler),
)

func provideSystemHandler(i do.Injector) (*systemhandler.Handler, error) {
	return systemhandler.NewHandler(), nil
}
