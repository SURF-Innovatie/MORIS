package di

import (
	"github.com/SURF-Innovatie/MORIS/internal/app/eventpolicy"
	"github.com/SURF-Innovatie/MORIS/internal/app/notification"
	organisationhierarchy "github.com/SURF-Innovatie/MORIS/internal/app/organisation/hierarchy"
	"github.com/SURF-Innovatie/MORIS/internal/domain/project/events/hydrator"
	eventpolicyadapter "github.com/SURF-Innovatie/MORIS/internal/infra/adapters/eventpolicy"
	eventpolicyrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/eventpolicy"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideEventPolicyService),
	do.Lazy(providePolicyEvaluator),
)

func provideEventPolicyService(i do.Injector) (eventpolicy.Service, error) {
	repo := do.MustInvoke[*eventpolicyrepo.EntRepo](i)
	orgHierarchySvc := do.MustInvoke[organisationhierarchy.Service](i)
	return eventpolicy.NewService(repo, orgHierarchySvc), nil
}

func providePolicyEvaluator(i do.Injector) (eventpolicy.Evaluator, error) {
	repo := do.MustInvoke[*eventpolicyrepo.EntRepo](i)
	orgHierarchySvc := do.MustInvoke[organisationhierarchy.Service](i)
	recipient := do.MustInvoke[*eventpolicyadapter.RecipientAdapter](i)
	notifSvc := do.MustInvoke[notification.Service](i)
	h := do.MustInvoke[*hydrator.Hydrator](i)
	return eventpolicy.NewEvaluator(repo, orgHierarchySvc, recipient, notifSvc, h), nil
}
