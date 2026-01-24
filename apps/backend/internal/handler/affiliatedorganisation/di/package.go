package di

import (
	"github.com/SURF-Innovatie/MORIS/internal/app/affiliatedorganisation"
	affiliatedorganisationhandler "github.com/SURF-Innovatie/MORIS/internal/handler/affiliatedorganisation"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(func(i do.Injector) (int, error) {
		do.Provide(i, func(i do.Injector) (*affiliatedorganisationhandler.Handler, error) {
			svc := do.MustInvoke[affiliatedorganisation.Service](i)
			return affiliatedorganisationhandler.NewHandler(svc), nil
		})
		return 0, nil
	}),
)
