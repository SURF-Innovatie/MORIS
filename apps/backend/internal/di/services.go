package di

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	excrossref "github.com/SURF-Innovatie/MORIS/external/crossref"
	exnwo "github.com/SURF-Innovatie/MORIS/external/nwo"
	exorcid "github.com/SURF-Innovatie/MORIS/external/orcid"
	exsurfconext "github.com/SURF-Innovatie/MORIS/external/surfconext"
	exzenodo "github.com/SURF-Innovatie/MORIS/external/zenodo"
	coreauth "github.com/SURF-Innovatie/MORIS/internal/app/auth"
	"github.com/SURF-Innovatie/MORIS/internal/app/crossref"
	"github.com/SURF-Innovatie/MORIS/internal/app/customfield"
	"github.com/SURF-Innovatie/MORIS/internal/app/doi"
	"github.com/SURF-Innovatie/MORIS/internal/app/errorlog"
	"github.com/SURF-Innovatie/MORIS/internal/app/eventpolicy"
	"github.com/SURF-Innovatie/MORIS/internal/app/notification"
	"github.com/SURF-Innovatie/MORIS/internal/app/nwo"
	"github.com/SURF-Innovatie/MORIS/internal/app/orcid"
	"github.com/SURF-Innovatie/MORIS/internal/app/organisation"
	organisationrbac "github.com/SURF-Innovatie/MORIS/internal/app/organisation/rbac"
	personsvc "github.com/SURF-Innovatie/MORIS/internal/app/person"
	"github.com/SURF-Innovatie/MORIS/internal/app/portfolio"
	"github.com/SURF-Innovatie/MORIS/internal/app/product"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/cachewarmup"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/command"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/load"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/queries"
	"github.com/SURF-Innovatie/MORIS/internal/app/projectrole"
	"github.com/SURF-Innovatie/MORIS/internal/app/surfconext"
	"github.com/SURF-Innovatie/MORIS/internal/app/user"
	"github.com/SURF-Innovatie/MORIS/internal/app/zenodo"
	"github.com/SURF-Innovatie/MORIS/internal/event"
	"github.com/SURF-Innovatie/MORIS/internal/infra/adapters/event_policy"
	"github.com/SURF-Innovatie/MORIS/internal/infra/auth"
	"github.com/SURF-Innovatie/MORIS/internal/infra/cache"
	"github.com/SURF-Innovatie/MORIS/internal/infra/env"
	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/entclient"
	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/enttx"
	eventpolicyrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/eventpolicy"
	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/eventstore"
	notificationrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/notification"
	organisationrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/organisation"
	organisationrbacrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/organisation_rbac"
	personrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/person"
	portfoliorepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/portfolio"
	productrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/product"
	projectmembershiprepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/project_membership"
	projectqueryrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/project_query"
	userrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/user"
	"github.com/samber/do/v2"
)

func providePersonService(i do.Injector) (personsvc.Service, error) {
	repo := do.MustInvoke[*personrepo.EntRepo](i)
	return personsvc.NewService(repo), nil
}

func provideUserService(i do.Injector) (user.Service, error) {
	userRepo := do.MustInvoke[*userrepo.EntRepo](i)
	personSvc := do.MustInvoke[personsvc.Service](i)
	es := do.MustInvoke[*eventstore.EntStore](i)
	membership := do.MustInvoke[*projectmembershiprepo.EntRepo](i)
	userCache := do.MustInvoke[cache.UserCache](i)
	return user.NewService(userRepo, personSvc, es, membership, userCache), nil
}

func provideAuthService(i do.Injector) (coreauth.Service, error) {
	cli := do.MustInvoke[*ent.Client](i)
	userSvc := do.MustInvoke[user.Service](i)
	personSvc := do.MustInvoke[personsvc.Service](i)
	return auth.NewJWTService(cli, userSvc, personSvc, env.Global.JWTSecret), nil
}

func provideCurrentUserProvider(i do.Injector) (coreauth.CurrentUserProvider, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return auth.NewCurrentUserProvider(cli), nil
}

func provideORCIDService(i do.Injector) (orcid.Service, error) {
	userRepo := do.MustInvoke[*userrepo.EntRepo](i)
	personRepo := do.MustInvoke[*personrepo.EntRepo](i)
	cli := do.MustInvoke[*exorcid.Client](i)
	return orcid.NewService(userRepo, personRepo, cli), nil
}

func provideSurfconextService(i do.Injector) (surfconext.Service, error) {
	cli := do.MustInvoke[*exsurfconext.Client](i)
	authSvc := do.MustInvoke[coreauth.Service](i)
	return surfconext.NewService(cli, authSvc), nil
}

func provideZenodoService(i do.Injector) (zenodo.Service, error) {
	userRepo := do.MustInvoke[*userrepo.EntRepo](i)
	cli := do.MustInvoke[*exzenodo.Client](i)
	return zenodo.NewService(userRepo, cli), nil
}

func provideCrossrefService(i do.Injector) (crossref.Service, error) {
	cli := do.MustInvoke[excrossref.Client](i)
	return crossref.NewService(cli), nil
}

func provideNWOService(i do.Injector) (nwo.Service, error) {
	cli := do.MustInvoke[exnwo.Client](i)
	return nwo.NewService(cli), nil
}

func provideDoiService(i do.Injector) (doi.Service, error) {
	return doi.NewService(), nil
}

func provideOrgRBACService(i do.Injector) (organisationrbac.Service, error) {
	repo := do.MustInvoke[*organisationrbacrepo.EntRepo](i)
	return organisationrbac.NewService(repo), nil
}

func provideProjectRoleService(i do.Injector) (projectrole.Service, error) {
	repo := do.MustInvoke[projectrole.Repository](i)
	orgSvc := do.MustInvoke[organisation.Service](i)
	return projectrole.NewService(repo, orgSvc), nil
}

func provideCustomFieldService(i do.Injector) (customfield.Service, error) {
	repo := do.MustInvoke[customfield.Repository](i)
	return customfield.NewService(repo), nil
}

func provideOrganisationService(i do.Injector) (organisation.Service, error) {
	orgRepo := do.MustInvoke[*organisationrepo.EntRepo](i)
	personSvc := do.MustInvoke[*personrepo.EntRepo](i)
	rbacSvc := do.MustInvoke[organisationrbac.Service](i)
	txManager := do.MustInvoke[*enttx.Manager](i)
	return organisation.NewService(orgRepo, personSvc, rbacSvc, txManager), nil
}

func provideProductService(i do.Injector) (product.Service, error) {
	repo := do.MustInvoke[*productrepo.EntRepo](i)
	return product.NewService(repo), nil
}

func providePortfolioService(i do.Injector) (portfolio.Service, error) {
	repo := do.MustInvoke[*portfoliorepo.EntRepo](i)
	return portfolio.NewService(repo), nil
}

func provideNotificationService(i do.Injector) (notification.Service, error) {
	repo := do.MustInvoke[*notificationrepo.EntRepo](i)
	return notification.NewService(repo), nil
}

func provideErrorLogService(i do.Injector) (errorlog.Service, error) {
	repo := do.MustInvoke[errorlog.Repository](i)
	return errorlog.NewService(repo), nil
}

func provideEventPolicyService(i do.Injector) (eventpolicy.Service, error) {
	repo := do.MustInvoke[*eventpolicyrepo.EntRepo](i)
	closure := do.MustInvoke[*event_policy.OrgClosureAdapter](i)
	return eventpolicy.NewService(repo, closure), nil
}

func providePolicyEvaluator(i do.Injector) (eventpolicy.Evaluator, error) {
	repo := do.MustInvoke[*eventpolicyrepo.EntRepo](i)
	closure := do.MustInvoke[*event_policy.OrgClosureAdapter](i)
	recipient := do.MustInvoke[*event_policy.RecipientAdapter](i)
	notifSvc := do.MustInvoke[notification.Service](i)
	return eventpolicy.NewEvaluator(repo, closure, recipient, notifSvc), nil
}

func provideProjectLoader(i do.Injector) (*load.Loader, error) {
	es := do.MustInvoke[*eventstore.EntStore](i)
	pc := do.MustInvoke[cache.ProjectCache](i)
	return load.New(es, pc), nil
}

func provideProjectQueryService(i do.Injector) (queries.Service, error) {
	es := do.MustInvoke[*eventstore.EntStore](i)
	ldr := do.MustInvoke[*load.Loader](i)
	repo := do.MustInvoke[*projectqueryrepo.EntRepo](i)
	roleRepo := do.MustInvoke[projectrole.Repository](i)
	curUser := do.MustInvoke[coreauth.CurrentUserProvider](i)
	userSvc := do.MustInvoke[user.Service](i)
	return queries.NewService(es, ldr, repo, roleRepo, curUser, userSvc), nil
}

func provideEntClientProvider(i do.Injector) (command.EntClientProvider, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return entclient.New(cli), nil
}

func provideTxManager(i do.Injector) (*enttx.Manager, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return enttx.NewManager(cli), nil
}

func provideProjectCommandService(i do.Injector) (command.Service, error) {
	es := do.MustInvoke[*eventstore.EntStore](i)
	evtSvc := do.MustInvoke[event.Service](i)
	pc := do.MustInvoke[cache.ProjectCache](i)
	ref := do.MustInvoke[cache.ProjectCacheRefresher](i)
	curUser := do.MustInvoke[coreauth.CurrentUserProvider](i)
	entProv := do.MustInvoke[command.EntClientProvider](i)
	roleSvc := do.MustInvoke[projectrole.Service](i)
	evaluator := do.MustInvoke[eventpolicy.Evaluator](i)
	orgSvc := do.MustInvoke[organisation.Service](i)
	rbacSvc := do.MustInvoke[organisationrbac.Service](i)
	return command.NewService(es, evtSvc, pc, ref, curUser, entProv, roleSvc, evaluator, orgSvc, rbacSvc), nil
}

func provideCacheWarmupService(i do.Injector) (cachewarmup.Service, error) {
	repo := do.MustInvoke[*projectqueryrepo.EntRepo](i)
	ldr := do.MustInvoke[*load.Loader](i)
	pc := do.MustInvoke[cache.ProjectCache](i)
	return cachewarmup.NewService(repo, ldr, pc), nil
}

func provideEventService(i do.Injector) (event.Service, error) {
	es := do.MustInvoke[*eventstore.EntStore](i)
	cli := do.MustInvoke[*ent.Client](i)
	notifSvc := do.MustInvoke[notification.Service](i)

	projHandler := do.MustInvoke[*event.ProjectEventNotificationHandler](i)
	approvalHandler := do.MustInvoke[*event.ApprovalRequestNotificationHandler](i)
	policyHandler := do.MustInvoke[*event.Handler](i)
	execHandler := do.MustInvoke[*event.PolicyExecutionHandler](i)

	notificationHandlers := []event.NotificationHandler{
		projHandler,
		approvalHandler,
		policyHandler,
		execHandler,
	}

	statusHandler := do.MustInvoke[*event.StatusUpdateNotificationHandler](i)
	cacheHandler := do.MustInvoke[*event.CacheRefreshHandler](i)

	statusChangeHandlers := []event.StatusChangeHandler{
		statusHandler.Handle,
		cacheHandler.Handle,
	}

	return event.NewService(es, cli, notifSvc, notificationHandlers, statusChangeHandlers), nil
}
