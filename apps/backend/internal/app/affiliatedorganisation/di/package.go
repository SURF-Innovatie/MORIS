package di

import (
	"github.com/SURF-Innovatie/MORIS/internal/app/affiliatedorganisation"
	"github.com/samber/do/v2"
)

var Package = do.Lazy(func(i do.Injector) (affiliatedorganisation.Service, error) {
	repo := do.MustInvoke[affiliatedorganisation.Repository](i)
	return affiliatedorganisation.NewService(repo), nil
})
