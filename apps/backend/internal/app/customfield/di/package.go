package di

import (
	"github.com/SURF-Innovatie/MORIS/internal/app/customfield"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideCustomFieldService),
)

func provideCustomFieldService(i do.Injector) (customfield.Service, error) {
	repo := do.MustInvoke[customfield.Repository](i)
	return customfield.NewService(repo), nil
}
