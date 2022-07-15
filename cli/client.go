package cli

import (
	"fmt"
	"os"

	"github.com/eventials/go-tus"
	"github.com/spf13/cobra"
	"github.com/wwqdrh/logger"
)

var clientFlag ClientCmdFlag

type ClientCmdFlag struct {
	host       string
	uploadfile string
}

var (
	// flag
	ClientCmd = &cobra.Command{
		Use:          "client",
		Short:        "start tusd client",
		Long:         "start tusd client, and upload a file",
		Example:      "...",
		SilenceUsage: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := ClientStart()
			if err != nil {
				logger.DefaultLogger.Error(err.Error())
			}
			return err
		},
	}
)

func init() {
	ClientCmd.Flags().StringVar(&clientFlag.host, "host", "", "目标ip")
	ClientCmd.Flags().StringVar(&clientFlag.uploadfile, "upload", "", "上传文件")

	if err := ClientCmd.MarkFlagRequired("host"); err != nil {
		logger.DefaultLogger.Error(err.Error())
	}
	if err := ClientCmd.MarkFlagRequired("upload"); err != nil {
		logger.DefaultLogger.Error(err.Error())
	}
}

func ClientStart() error {
	f, err := os.Open(clientFlag.uploadfile)
	if err != nil {
		return fmt.Errorf("打开目标文件失败: %w", err)
	}
	defer f.Close()

	client, err := tus.NewClient(clientFlag.host, nil)
	if err != nil {
		return fmt.Errorf("tus client初始化失败: %w", err)
	}
	upload, err := tus.NewUploadFromFile(f)
	if err != nil {
		return fmt.Errorf("tus client初始化文件上传失败: %w", err)
	}

	uploader, err := client.CreateUpload(upload)
	if err != nil {
		return fmt.Errorf("tus client初始化文件上传失败: %w", err)
	}
	err = uploader.Upload()
	if err != nil {
		return fmt.Errorf("tus client文件上传失败: %w", err)
	}

	return nil
}
