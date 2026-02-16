package di

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/app/affiliatedorganisation"
	persistence "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/affiliatedorganisation"
	"github.com/samber/do/v2"
)

var Package = do.Lazy(func(i do.Injector) (affiliatedorganisation.Repository, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return persistence.NewEntRepo(cli), nil
})
