package cli

import (
	"fmt"
	"net/http"

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
	ServerCmd.Flags().StringVar(&serverFlag.urlpath, "baseurl", "/files/", "url基础路径")
}

func ServerStart() error {
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
		return err
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
		logger.DefaultLogger.Error(err.Error())
		return err
	}
	return nil
}
