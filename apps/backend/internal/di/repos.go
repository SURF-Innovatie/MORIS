package di

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/app/customfield"
	"github.com/SURF-Innovatie/MORIS/internal/app/errorlog"
	"github.com/SURF-Innovatie/MORIS/internal/app/projectrole"
	customfieldrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/customfield"
	errorlogrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/error_log"
	eventpolicyrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/eventpolicy"
	notificationrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/notification"
	organisationrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/organisation"
	organisationrbacrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/organisation_rbac"
	personrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/person"
	portfoliorepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/portfolio"
	productrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/product"
	projectmembershiprepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/project_membership"
	projectqueryrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/project_query"
	projectrolerepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/projectrole"
	userrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/user"
	"github.com/samber/do/v2"
)

func providePersonRepo(i do.Injector) (*personrepo.EntRepo, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return personrepo.NewEntRepo(cli), nil
}

func provideUserRepo(i do.Injector) (*userrepo.EntRepo, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return userrepo.NewEntRepo(cli), nil
}

func provideMembershipRepo(i do.Injector) (*projectmembershiprepo.EntRepo, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return projectmembershiprepo.NewEntRepo(cli), nil
}

func provideOrgRepo(i do.Injector) (*organisationrepo.EntRepo, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return organisationrepo.NewEntRepo(cli), nil
}

func provideOrgRBACRepo(i do.Injector) (*organisationrbacrepo.EntRepo, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return organisationrbacrepo.NewEntRepo(cli), nil
}

func provideProjectRoleRepo(i do.Injector) (projectrole.Repository, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return projectrolerepo.NewRepository(cli), nil
}

func provideCustomFieldRepo(i do.Injector) (customfield.Repository, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return customfieldrepo.NewRepository(cli), nil
}

func provideProductRepo(i do.Injector) (*productrepo.EntRepo, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return productrepo.NewEntRepo(cli), nil
}

func providePortfolioRepo(i do.Injector) (*portfoliorepo.EntRepo, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return portfoliorepo.NewEntRepo(cli), nil
}

func provideNotificationRepo(i do.Injector) (*notificationrepo.EntRepo, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return notificationrepo.NewEntRepo(cli), nil
}

func provideErrorLogRepo(i do.Injector) (errorlog.Repository, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return errorlogrepo.NewRepository(cli), nil
}

func provideProjectQueryRepo(i do.Injector) (*projectqueryrepo.EntRepo, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return projectqueryrepo.NewEntRepo(cli), nil
}

func provideEventPolicyRepo(i do.Injector) (*eventpolicyrepo.EntRepo, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return eventpolicyrepo.NewEntRepository(cli), nil
}
