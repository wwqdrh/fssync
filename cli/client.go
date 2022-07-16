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
	}

	ClientDownloadCmd = &cobra.Command{
		Use:   "download",
		Short: "download",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := client.DownloadStart()
			if err != nil {
				logger.DefaultLogger.Error(err.Error())
			}
			return err
		},
	}

	ClientUploadCmd = &cobra.Command{
		Use:   "upload",
		Short: "upload",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := client.UploadStart()
			if err != nil {
				logger.DefaultLogger.Error(err.Error())
			}
			return err
		},
	}
)

func init() {
	ClientCmd.AddCommand(ClientUploadCmd)
	ClientCmd.AddCommand(ClientDownloadCmd)

	ClientUploadCmd.Flags().StringVar(&client.ClientUploadFlag.Host, "host", "", "目标地址 http://127.0.0.1:1080/files/")
	ClientUploadCmd.Flags().StringVar(&client.ClientUploadFlag.Uploadfile, "upload", "", "上传文件")
	ClientUploadCmd.Flags().StringVar(&client.ClientUploadFlag.SpecPath, "spec", "", "文件的分片信息")

	ClientDownloadCmd.Flags().StringVar(&client.ClientDownloadFlag.Host, "host", "", "目标地址 http://127.0.0.1:1080/files/")
	ClientDownloadCmd.Flags().StringVar(&client.ClientDownloadFlag.DownloadUrl, "url", "", "目标下载url http://127.0.0.1:1080/download")
	ClientDownloadCmd.Flags().StringVar(&client.ClientDownloadFlag.FileName, "filename", "", "目标下载文件名称 可从/download/list中查看")
	ClientDownloadCmd.Flags().StringVar(&client.ClientDownloadFlag.SpecPath, "spec", "", "文件的分片信息保存路径")

	if err := ClientUploadCmd.MarkFlagRequired("host"); err != nil {
		logger.DefaultLogger.Error(err.Error())
	}
	if err := ClientUploadCmd.MarkFlagRequired("upload"); err != nil {
		logger.DefaultLogger.Error(err.Error())
	}
}
