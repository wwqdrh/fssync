package client

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/wwqdrh/gokit/clitool"
	"github.com/wwqdrh/gokit/logger"
)

var (
	ClientDownloadFlag = ClientDownloadCmdFlag{
		DownloadUrl:  "http://127.0.0.1:1080",
		DownloadPath: ".",
		SpecPath:     "./tmp/spec",
		TempPath:     "./tmp/data",
		Interval:     10,
	}
	ClientUploadFlag ClientUploadCmdFlag
)

type ClientWebDavFlag struct {
}

type ClientDownloadCmdFlag struct {
	DownloadUrl  string `name:"host" desc:"服务端地址"`
	FileName     string `name:"file" desc:"需要下载的单个文件的名字"`
	DownloadPath string `name:"down" desc:"下载文件保存路径"`
	SpecPath     string `name:"spec" desc:"切片信息等保存路径"`
	TempPath     string `name:"temp" desc:"保存切片的临时目录"`             // 保存切片的临时目录
	DownAll      bool   `name:"all" desc:"是否下载所有的文件"`              // 下载所有的文件
	IsDel        bool   `name:"del" desc:"是否删除"`                   // 是否删除
	Watch        bool   `name:"watch" desc:"是监听更新"`                // 是否监听更新
	Interval     int    `name:"interval" alias:"i" desc:"更新数据的频率"` // 更新数据的频率
}

type ClientUploadCmdFlag struct {
	Host       string
	Uploadfile string
	SpecPath   string
}

func (c *ClientDownloadCmdFlag) Init() {
	var err error

	if c.DownloadPath, err = filepath.Abs(c.DownloadPath); err != nil {
		logger.DefaultLogger.Fatal(err.Error())
	}
	if c.SpecPath, err = filepath.Abs(c.SpecPath); err != nil {
		logger.DefaultLogger.Fatal(err.Error())
	}
	if c.TempPath, err = filepath.Abs(c.TempPath); err != nil {
		logger.DefaultLogger.Fatal(err.Error())
	}
}

func Command() *clitool.Command {
	cmd := &clitool.Command{
		Cmd: &cobra.Command{
			Use:   "client",
			Short: "client",
		},
	}
	cmd.Add(&clitool.Command{
		Cmd: &cobra.Command{
			Use:   "webdav",
			Short: "webdav",
			RunE: func(cmd *cobra.Command, args []string) error {
				return SyncWebdav()
			},
		},
	})
	cmd.Add(&clitool.Command{
		Cmd: &cobra.Command{
			Use:   "download",
			Short: "download",
			RunE: func(cmd *cobra.Command, args []string) error {
				return DownloadStart()
			},
		},
		Values: &ClientDownloadFlag,
	})
	return cmd
}
