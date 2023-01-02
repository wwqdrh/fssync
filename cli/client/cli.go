package client

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/wwqdrh/fssync/client"
	"github.com/wwqdrh/gokit/clitool"
)

func Command() *clitool.Command {
	cmd := &clitool.Command{
		Cmd: &cobra.Command{
			Use:   "client",
			Short: "client",
		},
	}
	cmd.Add(&clitool.Command{
		Cmd: &cobra.Command{
			Use:   "download",
			Short: "download",
			RunE: func(cmd *cobra.Command, args []string) error {
				return client.DownloadStart()
			},
		},
		Options: []clitool.OptionConfig{
			{
				Target:       "DownloadUrl",
				Name:         "host",
				DefaultValue: "http://127.0.0.1:1080",
				Description:  "服务端地址",
			},
			{
				Target:       "FileName",
				Name:         "file",
				DefaultValue: "",
				Description:  "需要下载的单个文件的名字",
			},
			{
				Target:       "DownloadPath",
				Name:         "down",
				DefaultValue: ".",
				Description:  "下载文件保存路径",
			},
			{
				Target:       "SpecPath",
				Name:         "spec",
				DefaultValue: "./tmpdata",
				Description:  "切片信息等保存路径",
			},
			{
				Target:       "TempPath",
				Name:         "temp",
				DefaultValue: "./tempdata",
				Description:  "保存切片的临时目录",
			},
			{
				Target:       "DownAll",
				Name:         "all",
				DefaultValue: false,
				Description:  "是否下载所有的文件",
			},
			{
				Target:       "IsDel",
				Name:         "del",
				DefaultValue: false,
				Description:  "是否删除",
			},
			{
				Target:       "Watch",
				Name:         "watch",
				DefaultValue: false,
				Description:  "是监听更新",
			},
			{
				Target:       "Interval",
				Name:         "interval",
				Alias:        "i",
				DefaultValue: 10,
				Description:  "更新数据的频率",
			},
		},
		Values: &client.ClientDownloadFlag,
	})
	cmd.Add(&clitool.Command{
		Cmd: &cobra.Command{
			Use:   "tui",
			Short: "tui",
			RunE: func(cmd *cobra.Command, args []string) error {
				return tea.NewProgram(NewClientView()).Start()
			},
		},
	})
	return cmd
}
