package di

import (
	coreauth "github.com/SURF-Innovatie/MORIS/internal/app/auth"
	bulkimportsvc "github.com/SURF-Innovatie/MORIS/internal/app/bulkimport"
	bulkimporthandler "github.com/SURF-Innovatie/MORIS/internal/handler/bulkimport"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideBulkImportHandler),
)

func provideBulkImportHandler(i do.Injector) (*bulkimporthandler.Handler, error) {
	svc := do.MustInvoke[bulkimportsvc.Service](i)
	curUser := do.MustInvoke[coreauth.CurrentUserProvider](i)
	return bulkimporthandler.NewHandler(svc, curUser), nil
}
