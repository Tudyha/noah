package service

import (
	"noah/client/app/service/logic/download"
	"noah/client/app/service/logic/explorer"
	"noah/client/app/service/logic/information"

	"github.com/samber/do/v2"
)

func Init(i do.Injector) {
	do.Provide(i, func(i do.Injector) (IInformationService, error) {
		return information.NewInformationService(i)
	})
	do.Provide(i, func(i do.Injector) (IFileExplorerService, error) {
		return explorer.NewFileService(i)
	})
	do.Provide(i, func(i do.Injector) (IDownloadService, error) {
		return download.NewDownloadService(i)
	})
}
