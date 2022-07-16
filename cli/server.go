// server
// 1、接收客户端文件上传: upload接口
// 2、支持客户端文件下载: listfile downloadfile
package cli

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/wwqdrh/fssync/server"
	"github.com/wwqdrh/logger"
)

var (
	// flag
	ServerCmd = &cobra.Command{
		Use:          "server",
		Short:        "start tusd server",
		Example:      "...",
		SilenceUsage: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithCancel(context.TODO())
			defer cancel()
			err := server.Start(ctx)
			if err != nil {
				logger.DefaultLogger.Error(err.Error())
			}
			return err
		},
	}
)

func init() {
	ServerCmd.Flags().StringVar(&server.ServerFlag.Port, "port", ":1080", "目标端口")
	ServerCmd.Flags().StringVar(&server.ServerFlag.Store, "store", "./stores", "保存路径")
	ServerCmd.Flags().StringVar(&server.ServerFlag.Urlpath, "baseurl", "/files/", "url基础路径")
	ServerCmd.Flags().StringVar(&server.ServerFlag.ExtraPath, "extra", "", "提供直接下载的文件夹路径")
}
