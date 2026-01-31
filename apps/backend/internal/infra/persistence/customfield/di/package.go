package di

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/app/customfield"
	customfieldrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/customfield"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideCustomFieldRepo),
)

func provideCustomFieldRepo(i do.Injector) (customfield.Repository, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return customfieldrepo.NewRepository(cli), nil
}
