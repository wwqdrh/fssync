package server

import (
	"path/filepath"

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
	Port          int
	Store         string
	Urlpath       string
	ExtraPath     string // 额外的直接下载的文件夹
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
}
