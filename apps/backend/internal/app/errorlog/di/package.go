package di

import (
	"github.com/SURF-Innovatie/MORIS/internal/app/errorlog"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideErrorLogService),
)

func provideErrorLogService(i do.Injector) (errorlog.Service, error) {
	repo := do.MustInvoke[errorlog.Repository](i)
	return errorlog.NewService(repo), nil
}
