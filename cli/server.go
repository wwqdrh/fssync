package cli

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/tus/tusd/pkg/filestore"
	tusd "github.com/tus/tusd/pkg/handler"
	"github.com/wwqdrh/logger"
)

var serverFlag ServerCmdFlag

type ServerCmdFlag struct {
	port    string
	store   string
	urlpath string
}

var (
	// flag
	ServerCmd = &cobra.Command{
		Use:          "server",
		Short:        "start tusd server",
		Example:      "...",
		SilenceUsage: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := ServerStart()
			if err != nil {
				logger.DefaultLogger.Error(err.Error())
			}
			return err
		},
	}
)

func init() {
	ServerCmd.Flags().StringVar(&serverFlag.port, "port", ":1080", "目标端口")
	ServerCmd.Flags().StringVar(&serverFlag.store, "store", "./stores", "保存路径")
	ServerCmd.Flags().StringVar(&serverFlag.urlpath, "baseurl", "/files", "url基础路径")
}

func ServerStart() error {
	if err := os.MkdirAll(serverFlag.store, 0o777); err != nil {
		return fmt.Errorf("创建保存路径失败: %w", err)
	}

	store := filestore.FileStore{
		Path: serverFlag.store,
	}
	composer := tusd.NewStoreComposer()
	store.UseIn(composer)

	handler, err := tusd.NewHandler(tusd.Config{
		BasePath:              serverFlag.urlpath,
		StoreComposer:         composer,
		NotifyCompleteUploads: true,
	})
	if err != nil {
		logger.DefaultLogger.Error(err.Error())
		return fmt.Errorf("创建tusd handler失败: %w", err)
	}

	go func() {
		for {
			event := <-handler.CompleteUploads
			logger.DefaultLogger.Info(fmt.Sprintf("Upload %s finished\n", event.Upload.ID))
		}
	}()

	http.Handle(serverFlag.urlpath, http.StripPrefix(serverFlag.urlpath, handler))
	logger.DefaultLogger.Info(serverFlag.port)
	err = http.ListenAndServe(serverFlag.port, nil)
	if err != nil {
		return fmt.Errorf("服务退出出错: %w", err)
	}
	return nil
}
