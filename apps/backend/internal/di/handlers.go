package di

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	coreauth "github.com/SURF-Innovatie/MORIS/internal/app/auth"
	"github.com/SURF-Innovatie/MORIS/internal/app/crossref"
	"github.com/SURF-Innovatie/MORIS/internal/app/customfield"
	"github.com/SURF-Innovatie/MORIS/internal/app/doi"
	"github.com/SURF-Innovatie/MORIS/internal/app/eventpolicy"
	"github.com/SURF-Innovatie/MORIS/internal/app/notification"
	"github.com/SURF-Innovatie/MORIS/internal/app/nwo"
	"github.com/SURF-Innovatie/MORIS/internal/app/orcid"
	"github.com/SURF-Innovatie/MORIS/internal/app/organisation"
	organisationrbac "github.com/SURF-Innovatie/MORIS/internal/app/organisation/rbac"
	organisationrole "github.com/SURF-Innovatie/MORIS/internal/app/organisation/role"
	personsvc "github.com/SURF-Innovatie/MORIS/internal/app/person"
	"github.com/SURF-Innovatie/MORIS/internal/app/portfolio"
	"github.com/SURF-Innovatie/MORIS/internal/app/product"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/command"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/queries"
	"github.com/SURF-Innovatie/MORIS/internal/app/projectrole"
	"github.com/SURF-Innovatie/MORIS/internal/app/surfconext"
	"github.com/SURF-Innovatie/MORIS/internal/app/user"
	"github.com/SURF-Innovatie/MORIS/internal/app/zenodo"
	"github.com/SURF-Innovatie/MORIS/internal/event"
	authhandler "github.com/SURF-Innovatie/MORIS/internal/handler/auth"
	crossrefhandler "github.com/SURF-Innovatie/MORIS/internal/handler/crossref"
	doihandler "github.com/SURF-Innovatie/MORIS/internal/handler/doi"
	eventhandler "github.com/SURF-Innovatie/MORIS/internal/handler/event"
	eventpolicyhandler "github.com/SURF-Innovatie/MORIS/internal/handler/eventpolicy"
	notificationhandler "github.com/SURF-Innovatie/MORIS/internal/handler/notification"
	nwohandler "github.com/SURF-Innovatie/MORIS/internal/handler/nwo"
	orcidhandler "github.com/SURF-Innovatie/MORIS/internal/handler/orcid"
	organisationhandler "github.com/SURF-Innovatie/MORIS/internal/handler/organisation"
	personhandler "github.com/SURF-Innovatie/MORIS/internal/handler/person"
	portfoliohandler "github.com/SURF-Innovatie/MORIS/internal/handler/portfolio"
	producthandler "github.com/SURF-Innovatie/MORIS/internal/handler/product"
	projecthandler "github.com/SURF-Innovatie/MORIS/internal/handler/project"
	commandhandler "github.com/SURF-Innovatie/MORIS/internal/handler/project/command"
	systemhandler "github.com/SURF-Innovatie/MORIS/internal/handler/system"
	userhandler "github.com/SURF-Innovatie/MORIS/internal/handler/user"
	zenodohandler "github.com/SURF-Innovatie/MORIS/internal/handler/zenodo"
	"github.com/SURF-Innovatie/MORIS/internal/infra/cache"
	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/eventstore"
	"github.com/samber/do/v2"
)

// Event Handler Providers

func provideProjectEventHandler(i do.Injector) (*event.ProjectEventNotificationHandler, error) {
	cli := do.MustInvoke[*ent.Client](i)
	es := do.MustInvoke[*eventstore.EntStore](i)
	return event.NewProjectEventHandler(cli, es), nil
}

func provideApprovalRequestHandler(i do.Injector) (*event.ApprovalRequestNotificationHandler, error) {
	cli := do.MustInvoke[*ent.Client](i)
	es := do.MustInvoke[*eventstore.EntStore](i)
	rbac := do.MustInvoke[organisationrbac.Service](i)
	return event.NewApprovalRequestHandler(cli, es, rbac), nil
}

func provideEventPolicyHandler(i do.Injector) (*event.Handler, error) {
	repo := do.MustInvoke[eventpolicy.Service](i)
	cli := do.MustInvoke[*ent.Client](i)
	return event.NewEventPolicyHandler(repo, cli), nil
}

func providePolicyExecutionHandler(i do.Injector) (*event.PolicyExecutionHandler, error) {
	evaluator := do.MustInvoke[eventpolicy.Evaluator](i)
	projSvc := do.MustInvoke[queries.Service](i)
	return event.NewPolicyExecutionHandler(evaluator, projSvc), nil
}

func provideStatusUpdateHandler(i do.Injector) (*event.StatusUpdateNotificationHandler, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return event.NewStatusUpdateHandler(cli), nil
}

func provideCacheRefreshHandler(i do.Injector) (*event.CacheRefreshHandler, error) {
	refresher := do.MustInvoke[cache.ProjectCacheRefresher](i)
	return event.NewCacheRefreshHandler(refresher), nil
}

// HTTP Handler Providers

func providePersonHandler(i do.Injector) (*personhandler.Handler, error) {
	svc := do.MustInvoke[personsvc.Service](i)
	return personhandler.NewHandler(svc), nil
}

func provideUserHandler(i do.Injector) (*userhandler.Handler, error) {
	userSvc := do.MustInvoke[user.Service](i)
	projSvc := do.MustInvoke[queries.Service](i)
	return userhandler.NewHandler(userSvc, projSvc), nil
}

func provideAuthHandler(i do.Injector) (*authhandler.Handler, error) {
	userSvc := do.MustInvoke[user.Service](i)
	authSvc := do.MustInvoke[coreauth.Service](i)
	orcidSvc := do.MustInvoke[orcid.Service](i)
	surfSvc := do.MustInvoke[surfconext.Service](i)
	return authhandler.NewHandler(userSvc, authSvc, orcidSvc, surfSvc), nil
}

func provideORCIDHandler(i do.Injector) (*orcidhandler.Handler, error) {
	svc := do.MustInvoke[orcid.Service](i)
	return orcidhandler.NewHandler(svc), nil
}

func provideZenodoHandler(i do.Injector) (*zenodohandler.Handler, error) {
	svc := do.MustInvoke[zenodo.Service](i)
	curUser := do.MustInvoke[coreauth.CurrentUserProvider](i)
	return zenodohandler.NewHandler(svc, curUser), nil
}

func provideCrossrefHandler(i do.Injector) (*crossrefhandler.Handler, error) {
	svc := do.MustInvoke[crossref.Service](i)
	return crossrefhandler.NewHandler(svc), nil
}

func provideNWOHandler(i do.Injector) (*nwohandler.Handler, error) {
	svc := do.MustInvoke[nwo.Service](i)
	return nwohandler.NewHandler(svc), nil
}

func provideDoiHandler(i do.Injector) (*doihandler.Handler, error) {
	svc := do.MustInvoke[doi.Service](i)
	return doihandler.NewHandler(svc), nil
}

func provideOrgRBACHandler(i do.Injector) (*organisationhandler.RBACHandler, error) {
	svc := do.MustInvoke[organisationrbac.Service](i)
	return organisationhandler.NewRBACHandler(svc), nil
}

func provideOrgRoleHandler(i do.Injector) (*organisationhandler.RoleHandler, error) {
	roleSvc := do.MustInvoke[organisationrole.Service](i)
	rbacSvc := do.MustInvoke[organisationrbac.Service](i)
	return organisationhandler.NewRoleHandler(roleSvc, rbacSvc), nil
}

func provideOrganisationHandler(i do.Injector) (*organisationhandler.Handler, error) {
	orgSvc := do.MustInvoke[organisation.Service](i)
	rbacSvc := do.MustInvoke[organisationrbac.Service](i)
	roleSvc := do.MustInvoke[projectrole.Service](i)
	cfSvc := do.MustInvoke[customfield.Service](i)
	return organisationhandler.NewHandler(orgSvc, rbacSvc, roleSvc, cfSvc), nil
}

func provideProductHandler(i do.Injector) (*producthandler.Handler, error) {
	svc := do.MustInvoke[product.Service](i)
	curUser := do.MustInvoke[coreauth.CurrentUserProvider](i)
	return producthandler.NewHandler(svc, curUser), nil
}

func providePortfolioHandler(i do.Injector) (*portfoliohandler.Handler, error) {
	svc := do.MustInvoke[portfolio.Service](i)
	curUser := do.MustInvoke[coreauth.CurrentUserProvider](i)
	return portfoliohandler.NewHandler(svc, curUser), nil
}

func provideNotificationHandler(i do.Injector) (*notificationhandler.Handler, error) {
	svc := do.MustInvoke[notification.Service](i)
	return notificationhandler.NewHandler(svc), nil
}

func provideEventHandler(i do.Injector) (*eventhandler.Handler, error) {
	evtSvc := do.MustInvoke[event.Service](i)
	projSvc := do.MustInvoke[queries.Service](i)
	userSvc := do.MustInvoke[user.Service](i)
	cli := do.MustInvoke[*ent.Client](i)
	return eventhandler.NewHandler(evtSvc, projSvc, userSvc, cli), nil
}

func provideEventPolicyHTTPHandler(i do.Injector) (*eventpolicyhandler.Handler, error) {
	svc := do.MustInvoke[eventpolicy.Service](i)
	return eventpolicyhandler.NewHandler(svc), nil
}

func provideProjectHandler(i do.Injector) (*projecthandler.Handler, error) {
	svc := do.MustInvoke[queries.Service](i)
	cfSvc := do.MustInvoke[customfield.Service](i)
	return projecthandler.NewHandler(svc, cfSvc), nil
}

func provideProjectCommandHandler(i do.Injector) (*commandhandler.Handler, error) {
	svc := do.MustInvoke[command.Service](i)
	return commandhandler.NewHandler(svc), nil
}

func provideSystemHandler(i do.Injector) (*systemhandler.Handler, error) {
	return systemhandler.NewHandler(), nil
}
