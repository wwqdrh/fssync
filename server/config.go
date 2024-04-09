package server

import (
	"context"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/wwqdrh/gokit/clitool"
	"github.com/wwqdrh/gokit/logger"
)

var ServerFlag = serverCmdFlag{
	Port:          1080,
	ExtraTruncate: 1 * 1024 * 1024,
	Urlpath:       "/files",
	ExtraPath:     ".",
	Store:         "./upload",
}

type serverCmdFlag struct {
	Port          int    `name:"port" alias:"p" desc:"端口号"`
	Store         string `name:"store" alias:"s" desc:"上传文件保存路径"`
	Urlpath       string
	ExtraPath     string `name:"download" alias:"d" desc:"提供下载功能的文件夹"` // 额外的直接下载的文件夹
	ExtraTruncate int64  // 额外的直接下载的文件夹 分片的大小
}

func (c *serverCmdFlag) Init() {
	var err error

	if c.Urlpath, err = filepath.Abs(c.Urlpath); err != nil {
		logger.DefaultLogger.Fatal(err.Error())
	}
	if c.ExtraPath, err = filepath.Abs(c.ExtraPath); err != nil {
		logger.DefaultLogger.Fatal(err.Error())
	}
	if c.Store, err = filepath.Abs(c.Store); err != nil {
		logger.DefaultLogger.Fatal(err.Error())
	}

	logger.DefaultLogger.Info(c.Store)
	logger.DefaultLogger.Info(c.ExtraPath)
}

func Command() *clitool.Command {
	return &clitool.Command{
		Cmd: &cobra.Command{
			Use:   "server",
			Short: "server",
			RunE: func(cmd *cobra.Command, args []string) error {
				ctx, cancel := context.WithCancel(context.TODO())
				defer cancel()
				return NewFileManager().Start(ctx)
			},
		},
		Values: &ServerFlag,
	}
}
