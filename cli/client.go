package cli

import (
	"github.com/spf13/cobra"
	"github.com/wwqdrh/fssync/client"
	"github.com/wwqdrh/logger"
)

var (
	// flag
	ClientCmd = &cobra.Command{
		Use:          "client",
		Short:        "start tusd client",
		Long:         "start tusd client, and upload a file",
		Example:      "...",
		SilenceUsage: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := client.Start()
			if err != nil {
				logger.DefaultLogger.Error(err.Error())
			}
			return err
		},
	}
)

func init() {
	ClientCmd.Flags().StringVar(&client.ClientFlag.Host, "host", "", "目标地址 http://127.0.0.1:1080/files/")
	ClientCmd.Flags().StringVar(&client.ClientFlag.Uploadfile, "upload", "", "上传文件")
	ClientCmd.Flags().StringVar(&client.ClientFlag.SpecPath, "spec", "", "文件的分片信息")

	if err := ClientCmd.MarkFlagRequired("host"); err != nil {
		logger.DefaultLogger.Error(err.Error())
	}
	if err := ClientCmd.MarkFlagRequired("upload"); err != nil {
		logger.DefaultLogger.Error(err.Error())
	}
}
