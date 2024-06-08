package client

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/wwqdrh/fssync/driver"
	"github.com/wwqdrh/gokit/logger"
	"github.com/wwqdrh/gokit/ostool/fileindex"
)

func registerFileUpdate(d driver.IDriver, ctx context.Context) {
	tree := fileindex.NewFileInfoTree(ClientWebDavFlag.Work, ClientWebDavFlag.Interval, ClientWebDavFlag.Ignores)
	tree.SetDefaultTimeLines(d.GetLastTimelineMap())
	tree.SetOnFileInfoUpdate(func(fi fileindex.FileIndex) {
		if strings.HasPrefix(fi.BaseName, ".") || strings.HasPrefix(fi.BaseName, "_") {
			// 隐藏文件不进行上传
			return
		}

		realPath, err := filepath.Rel(ClientWebDavFlag.Work, fi.Path)
		if err != nil {
			logger.DefaultLogger.Error(err.Error())
			return
		}
		logger.DefaultLogger.Infox("update file %s to %s", nil, fi.Path, realPath)
		go d.Update(fi.Path, realPath)
	})
	// 启动更新
	tree.Start()
	defer tree.Stop()
	<-ctx.Done()
}

func SyncWebdav() error {
	d, err := driver.LoadDriver(ClientWebDavFlag.CfgFile, ClientWebDavFlag.Name)
	if err != nil {
		return err
	}
	ctx, ctxCancel := context.WithCancel(context.TODO())
	go registerFileUpdate(d, ctx)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	<-sigChan
	ctxCancel()

	return nil
}
