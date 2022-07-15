package server

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/tus/tusd/pkg/filestore"
	tusd "github.com/tus/tusd/pkg/handler"
	"github.com/wwqdrh/logger"
)

func Start(ctx context.Context) error {
	if err := os.MkdirAll(ServerFlag.Store, 0o777); err != nil {
		return fmt.Errorf("创建保存路径失败: %w", err)
	}

	store := filestore.FileStore{
		Path: ServerFlag.Store,
	}
	composer := tusd.NewStoreComposer()
	store.UseIn(composer)

	handler, err := tusd.NewHandler(tusd.Config{
		BasePath:              ServerFlag.Urlpath,
		StoreComposer:         composer,
		NotifyCompleteUploads: true,
	})
	if err != nil {
		logger.DefaultLogger.Error(err.Error())
		return fmt.Errorf("创建tusd handler失败: %w", err)
	}

	go func() {
		for {
			select {
			case event := <-handler.CompleteUploads:
				logger.DefaultLogger.Info(fmt.Sprintf("Upload %s finished\n", event.Upload.ID))
			case <-ctx.Done():
				return
			}
		}
	}()

	http.Handle(ServerFlag.Urlpath, http.StripPrefix(ServerFlag.Urlpath, handler))

	go func() {
		logger.DefaultLogger.Info(ServerFlag.Port)
		err = http.ListenAndServe(ServerFlag.Port, nil)
		if err != nil {
			logger.DefaultLogger.Error(fmt.Sprintf("服务退出出错: %s", err.Error()))
		}
	}()

	<-ctx.Done()

	return nil
}
