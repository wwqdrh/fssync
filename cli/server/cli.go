package server

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/wwqdrh/fssync/server"
	"github.com/wwqdrh/gokit/clitool"
)

func Command() *clitool.Command {
	return &clitool.Command{
		Cmd: &cobra.Command{
			Use:   "server",
			Short: "server",
			RunE: func(cmd *cobra.Command, args []string) error {
				ctx, cancel := context.WithCancel(context.TODO())
				defer cancel()
				return server.Start(ctx)
			},
		},
		Options: []clitool.OptionConfig{
			{
				Target:       "Port",
				Name:         "port",
				Alias:        "p",
				DefaultValue: 1080,
				Description:  "端口号",
			},
			{
				Target:       "Store",
				Name:         "store",
				Alias:        "s",
				DefaultValue: "./stores",
				Description:  "上传文件保存路径: (./stores)",
			},
			{
				Target:       "ExtraPath",
				Name:         "download",
				Alias:        "d",
				DefaultValue: "",
				Description:  "提供下载功能的文件夹: ('')",
			},
		},
		Values: &server.ServerFlag,
	}
}
