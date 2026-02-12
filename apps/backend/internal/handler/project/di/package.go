package di

import (
	coreauth "github.com/SURF-Innovatie/MORIS/internal/app/auth"
	"github.com/SURF-Innovatie/MORIS/internal/app/customfield"
	bulkimportsvc "github.com/SURF-Innovatie/MORIS/internal/app/project/bulkimport"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/command"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/queries"
	projecthandler "github.com/SURF-Innovatie/MORIS/internal/handler/project"
	bulkimporthandler "github.com/SURF-Innovatie/MORIS/internal/handler/project/bulkimport"
	commandhandler "github.com/SURF-Innovatie/MORIS/internal/handler/project/command"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideProjectHandler),
	do.Lazy(provideProjectCommandHandler),
	do.Lazy(provideBulkImportHandler),
)

func provideProjectHandler(i do.Injector) (*projecthandler.Handler, error) {
	svc := do.MustInvoke[queries.Service](i)
	cfSvc := do.MustInvoke[customfield.Service](i)
	return projecthandler.NewHandler(svc, cfSvc), nil
}

func provideProjectCommandHandler(i do.Injector) (*commandhandler.Handler, error) {
	svc := do.MustInvoke[command.Service](i)
	return commandhandler.NewHandler(svc), nil
}

func provideBulkImportHandler(i do.Injector) (*bulkimporthandler.Handler, error) {
	svc := do.MustInvoke[bulkimportsvc.Service](i)
	curUser := do.MustInvoke[coreauth.CurrentUserProvider](i)
	return bulkimporthandler.NewHandler(svc, curUser), nil
}
