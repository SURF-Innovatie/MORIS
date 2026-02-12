package wire

import (
	crossrefclientdi "github.com/SURF-Innovatie/MORIS/external/crossref/di"
	doiclientdi "github.com/SURF-Innovatie/MORIS/external/doi/di"
	nwoclientdi "github.com/SURF-Innovatie/MORIS/external/nwo/di"
	orcidoclientdi "github.com/SURF-Innovatie/MORIS/external/orcid/di"
	raidclientdi "github.com/SURF-Innovatie/MORIS/external/raid/di"
	surfconextclientdi "github.com/SURF-Innovatie/MORIS/external/surfconext/di"
	zenodoclientdi "github.com/SURF-Innovatie/MORIS/external/zenodo/di"
	adapterinternaldi "github.com/SURF-Innovatie/MORIS/internal/adapter/di"
	authappdi "github.com/SURF-Innovatie/MORIS/internal/app/auth/di"
	bulkimportappdi "github.com/SURF-Innovatie/MORIS/internal/app/bulkimport/di"
	crossrefappdi "github.com/SURF-Innovatie/MORIS/internal/app/crossref/di"
	customfieldappdi "github.com/SURF-Innovatie/MORIS/internal/app/customfield/di"
	doiappdi "github.com/SURF-Innovatie/MORIS/internal/app/doi/di"
	errorlogappdi "github.com/SURF-Innovatie/MORIS/internal/app/errorlog/di"
	eventappdi "github.com/SURF-Innovatie/MORIS/internal/app/event/di"
	eventpolicyappdi "github.com/SURF-Innovatie/MORIS/internal/app/eventpolicy/di"
	notificaionappdi "github.com/SURF-Innovatie/MORIS/internal/app/notification/di"
	nwoappdi "github.com/SURF-Innovatie/MORIS/internal/app/nwo/di"
	orcidappdi "github.com/SURF-Innovatie/MORIS/internal/app/orcid/di"
	organisationappdi "github.com/SURF-Innovatie/MORIS/internal/app/organisation/di"
	personappdi "github.com/SURF-Innovatie/MORIS/internal/app/person/di"
	portfolioappdi "github.com/SURF-Innovatie/MORIS/internal/app/portfolio/di"
	productappdi "github.com/SURF-Innovatie/MORIS/internal/app/product/di"
	projectappdi "github.com/SURF-Innovatie/MORIS/internal/app/project/di"
	surfconextappdi "github.com/SURF-Innovatie/MORIS/internal/app/surfconext/di"
	userappdi "github.com/SURF-Innovatie/MORIS/internal/app/user/di"
	zenodoappdi "github.com/SURF-Innovatie/MORIS/internal/app/zenodo/di"
	adapterhandlerdi "github.com/SURF-Innovatie/MORIS/internal/handler/adapter/di"
	authhandlerdi "github.com/SURF-Innovatie/MORIS/internal/handler/auth/di"
	bulkimporthandlerdi "github.com/SURF-Innovatie/MORIS/internal/handler/bulkimport/di"
	crossrefhandlerdi "github.com/SURF-Innovatie/MORIS/internal/handler/crossref/di"
	doihandlerdi "github.com/SURF-Innovatie/MORIS/internal/handler/doi/di"
	eventhandlerdi "github.com/SURF-Innovatie/MORIS/internal/handler/event/di"
	eventpolicyhandlerdi "github.com/SURF-Innovatie/MORIS/internal/handler/eventpolicy/di"
	notificationhandlerdi "github.com/SURF-Innovatie/MORIS/internal/handler/notification/di"
	nwohandlerdi "github.com/SURF-Innovatie/MORIS/internal/handler/nwo/di"
	orcidhandlerdi "github.com/SURF-Innovatie/MORIS/internal/handler/orcid/di"
	organisationhandlerdi "github.com/SURF-Innovatie/MORIS/internal/handler/organisation/di"
	personhandlerdi "github.com/SURF-Innovatie/MORIS/internal/handler/person/di"
	portfoliohandlerdi "github.com/SURF-Innovatie/MORIS/internal/handler/portfolio/di"
	producthandlerdi "github.com/SURF-Innovatie/MORIS/internal/handler/product/di"
	projecthandlerdi "github.com/SURF-Innovatie/MORIS/internal/handler/project/di"
	systemhandlerdi "github.com/SURF-Innovatie/MORIS/internal/handler/system/di"
	userhandlerdi "github.com/SURF-Innovatie/MORIS/internal/handler/user/di"
	zenodohandlerdi "github.com/SURF-Innovatie/MORIS/internal/handler/zenodo/di"
	recipientadapterinfradi "github.com/SURF-Innovatie/MORIS/internal/infra/adapters/eventpolicy/di"
	infradi "github.com/SURF-Innovatie/MORIS/internal/infra/di"
	eventpublisherinfradi "github.com/SURF-Innovatie/MORIS/internal/infra/eventdispatch/di"
	eventinfrahandlerdi "github.com/SURF-Innovatie/MORIS/internal/infra/handlers/events/di"
	identityinfradi "github.com/SURF-Innovatie/MORIS/internal/infra/identity/di"
	authrepodi "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/auth/di"
	customfieldrepodi "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/customfield/di"
	errorlogrepodi "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/errorlog/di"
	eventrepodi "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/event/di"
	eventpolicyrepodi "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/eventpolicy/di"
	notificationrepodi "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/notification/di"
	organisationrepodi "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/organisation/di"
	personrepodi "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/person/di"
	portfolierepodi "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/portfolio/di"
	productrepodi "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/product/di"
	projectrepodi "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/project/di"
	userrepodi "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/user/di"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	infradi.Package,
	adapterinternaldi.Package,

	adapterhandlerdi.Package,

	authappdi.Package,
	authhandlerdi.Package,
	authrepodi.Package,

	bulkimportappdi.Package,
	bulkimporthandlerdi.Package,

	crossrefappdi.Package,
	crossrefhandlerdi.Package,
	crossrefclientdi.Package,

	customfieldappdi.Package,
	customfieldrepodi.Package,

	doiappdi.Package,
	doihandlerdi.Package,
	doiclientdi.Package,

	errorlogappdi.Package,
	errorlogrepodi.Package,

	eventappdi.Package,
	eventrepodi.Package,
	eventhandlerdi.Package,
	eventinfrahandlerdi.Package,

	eventpolicyappdi.Package,
	eventpolicyrepodi.Package,
	eventpolicyhandlerdi.Package,

	eventpublisherinfradi.Package,

	identityinfradi.Package,

	notificaionappdi.Package,
	notificationrepodi.Package,
	notificationhandlerdi.Package,

	nwoappdi.Package,
	nwohandlerdi.Package,
	nwoclientdi.Package,

	orcidappdi.Package,
	orcidhandlerdi.Package,
	orcidoclientdi.Package,

	organisationappdi.Package,
	organisationrepodi.Package,
	organisationhandlerdi.Package,

	personrepodi.Package,
	personappdi.Package,
	personhandlerdi.Package,

	portfolioappdi.Package,
	portfolierepodi.Package,
	portfoliohandlerdi.Package,

	productappdi.Package,
	productrepodi.Package,
	producthandlerdi.Package,

	projectappdi.Package,
	projectrepodi.Package,
	projecthandlerdi.Package,

	raidclientdi.Package,

	recipientadapterinfradi.Package,

	systemhandlerdi.Package,

	surfconextappdi.Package,
	surfconextclientdi.Package,

	userappdi.Package,
	userrepodi.Package,
	userhandlerdi.Package,

	zenodoappdi.Package,
	zenodohandlerdi.Package,
	zenodoclientdi.Package,
)
