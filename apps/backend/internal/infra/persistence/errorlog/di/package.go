package di

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/app/errorlog"
	errorlogrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/errorlog"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideErrorLogRepo),
)

func provideErrorLogRepo(i do.Injector) (errorlog.Repository, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return errorlogrepo.NewRepository(cli), nil
}
