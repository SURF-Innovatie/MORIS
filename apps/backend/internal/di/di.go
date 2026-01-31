// Package di provides dependency injection using samber/do v2.
// All services and handlers are registered as lazy providers to optimize startup.
package di

import "github.com/samber/do/v2"

// Package combines all DI providers
var Package = do.Package(
	// Infrastructure
	do.Lazy(provideEntClient),
	do.Lazy(provideRedisClient),
	do.Lazy(ProvideEventRepo),
	do.Lazy(provideProjectCache),
	do.Lazy(provideUserCache),
	do.Lazy(provideCacheRefresher),

	// Repositories
	do.Lazy(providePersonRepo),
	do.Lazy(provideUserRepo),
	do.Lazy(provideMembershipRepo),
	do.Lazy(provideOrgRepo),
	do.Lazy(provideOrgRBACRepo),
	do.Lazy(provideOrgHierarchyRepo),
	do.Lazy(provideOrgRoleRepo),
	do.Lazy(provideProjectRoleRepo),
	do.Lazy(provideCustomFieldRepo),
	do.Lazy(provideProductRepo),
	do.Lazy(providePortfolioRepo),
	do.Lazy(provideNotificationRepo),
	do.Lazy(provideErrorLogRepo),
	do.Lazy(provideProjectRepo),
	do.Lazy(provideEventPolicyRepo),

	// Event Publisher
	do.Lazy(provideEventPublisher),

	// External Clients
	do.Lazy(provideORCIDClient),
	do.Lazy(provideSurfconextClient),
	do.Lazy(provideZenodoClient),
	do.Lazy(provideCrossrefClient),
	do.Lazy(provideRAiDClient),
	do.Lazy(provideNWOClient),

	// Adapters
	do.Lazy(provideRecipientAdapter),
	do.Lazy(provideAdapterRegistry),
	do.Lazy(provideAdapterHandler),

	// Services
	do.Lazy(providePersonService),
	do.Lazy(provideUserService),
	do.Lazy(provideAuthService),
	do.Lazy(provideCurrentUserProvider),
	do.Lazy(provideORCIDService),
	do.Lazy(provideSurfconextService),
	do.Lazy(provideZenodoService),
	do.Lazy(provideCrossrefService),
	do.Lazy(provideNWOService),
	do.Lazy(provideDoiService),
	do.Lazy(provideOrgRBACService),
	do.Lazy(provideOrgRoleService),
	do.Lazy(provideOrgHierarchyService),
	do.Lazy(provideProjectRoleService),
	do.Lazy(provideCustomFieldService),
	do.Lazy(provideOrganisationService),
	do.Lazy(provideProductService),
	do.Lazy(providePortfolioService),
	do.Lazy(provideNotificationService),
	do.Lazy(provideErrorLogService),
	do.Lazy(provideEventPolicyService),
	do.Lazy(providePolicyEvaluator),
	do.Lazy(provideProjectLoader),
	do.Lazy(provideProjectQueryService),
	do.Lazy(provideEntClientProvider),
	do.Lazy(provideTxManager),

	// Event handlers
	do.Lazy(provideProjectEventHandler),
	do.Lazy(provideApprovalRequestHandler),
	do.Lazy(provideEventPolicyHandler),
	do.Lazy(providePolicyExecutionHandler),
	do.Lazy(provideStatusUpdateHandler),
	do.Lazy(provideCacheRefreshHandler),

	// Event Service (depends on handlers)
	do.Lazy(provideEventService),

	// Command service (depends on event service)
	do.Lazy(provideProjectCommandService),
	do.Lazy(provideCacheWarmupService),

	// HTTP Handlers
	do.Lazy(providePersonHandler),
	do.Lazy(provideUserHandler),
	do.Lazy(provideAuthHandler),
	do.Lazy(provideORCIDHandler),
	do.Lazy(provideZenodoHandler),
	do.Lazy(provideCrossrefHandler),
	do.Lazy(provideNWOHandler),
	do.Lazy(provideDoiHandler),
	do.Lazy(provideOrgRBACHandler),
	do.Lazy(provideOrgRoleHandler),
	do.Lazy(provideOrganisationHandler),
	do.Lazy(provideProductHandler),
	do.Lazy(providePortfolioHandler),
	do.Lazy(provideNotificationHandler),
	do.Lazy(provideEventHandler),
	do.Lazy(provideEventPolicyHTTPHandler),
	do.Lazy(provideProjectHandler),
	do.Lazy(provideProjectCommandHandler),
	do.Lazy(provideSystemHandler),
)
