package client

import (
	"path/filepath"

	"github.com/wwqdrh/gokit/logger"
)

var ClientDownloadFlag = clientDownloadCmdFlag{
	DownloadPath: ".",
	SpecPath:     "./tmp/spec",
	TempPath:     "./tmp/data",
	Interval:     10,
}
var ClientUploadFlag clientUploadCmdFlag

type clientDownloadCmdFlag struct {
	DownloadUrl string
	FileName    string

	DownloadPath string
	SpecPath     string
	TempPath     string // 保存切片的临时目录
	DownAll      bool   // 下载所有的文件
	IsDel        bool   // 是否删除
	Watch        bool   // 是否监听更新
	Interval     int    // 更新数据的频率
}

type clientUploadCmdFlag struct {
	Host       string
	Uploadfile string
	SpecPath   string
}

func (c *clientDownloadCmdFlag) Init() {
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
