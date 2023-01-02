package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/tus/tusd/pkg/filestore"
	tusd "github.com/tus/tusd/pkg/handler"
	"github.com/wwqdrh/fssync/pkg/protocol"
	"github.com/wwqdrh/gokit/logger"
	"github.com/wwqdrh/gokit/ostool"
	"github.com/wwqdrh/gokit/ostool/fileindex"
)

type FileManager struct {
	updateFile *sync.Map
}

func NewFileManager() *FileManager {
	return &FileManager{
		updateFile: &sync.Map{},
	}
}

func (f *FileManager) Start(ctx context.Context) error {
	ServerFlag.Init()

	if err := os.MkdirAll(ServerFlag.Store, 0o777); err != nil {
		return fmt.Errorf("创建保存路径失败: %w", err)
	}

	store := filestore.FileStore{
		Path: ServerFlag.Store,
	}
	composer := tusd.NewStoreComposer()
	store.UseIn(composer)

	handler, err := tusd.NewHandler(tusd.Config{
		BasePath:              ServerFlag.Urlpath,
		StoreComposer:         composer,
		NotifyCompleteUploads: true,
	})
	if err != nil {
		logger.DefaultLogger.Error(err.Error())
		return fmt.Errorf("创建tusd handler失败: %w", err)
	}

	go func() {
		for {
			select {
			case event := <-handler.CompleteUploads:
				logger.DefaultLogger.Debug(fmt.Sprintf("Upload %s finished\n", event.Upload.ID))
			case <-ctx.Done():
				return
			}
		}
	}()

	http.Handle(ServerFlag.Urlpath, http.StripPrefix(ServerFlag.Urlpath, handler))
	f.registerAPI()

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	if err := f.watchFileModify(ctx); err != nil {
		logger.DefaultLogger.Warn(err.Error())
	}

	go func() {
		logger.DefaultLogger.Info(fmt.Sprintf("start server on: %d", ServerFlag.Port))
		err = http.ListenAndServe(fmt.Sprintf(":%d", ServerFlag.Port), nil)
		if err != nil {
			logger.DefaultLogger.Error(fmt.Sprintf("服务退出出错: %s", err.Error()))
		}
	}()

	n := make(chan os.Signal, 1)
	signal.Notify(n, syscall.SIGTERM, os.Interrupt)
	select {
	case <-ctx.Done():
		return nil
	case <-n:
		return nil
	}
}

func (f *FileManager) watchFileModify(ctx context.Context) error {
	return ostool.RegisterNotify(ctx, ServerFlag.ExtraPath, 1*time.Hour, func(e fsnotify.Event) {
		if abs, err := filepath.Abs(e.Name); err != nil {
			logger.DefaultLogger.Warn(err.Error())
		} else {
			f.updateFile.Store(strings.TrimPrefix(abs, ServerFlag.ExtraPath), struct{}{})
		}
	})
}

func (f *FileManager) registerAPI() {
	http.HandleFunc(protocol.PDownloadUpdate.ServerUrl(), f.downloadUpdateList)
	http.HandleFunc(protocol.PDownloadList.ServerUrl(), f.downloadList)
	http.HandleFunc(protocol.PDownloadSpec.ServerUrl(), f.downloadSpec)
	http.HandleFunc(protocol.PDownloadMd5.ServerUrl(), f.downloadMd5)
	http.HandleFunc(protocol.PDownloadTrucate.ServerUrl(), f.downloadTruncate)
	http.HandleFunc(protocol.PDownloadDelete.ServerUrl(), f.downloadDelete)
}

func (f *FileManager) downloadUpdateList(w http.ResponseWriter, r *http.Request) {
	if ServerFlag.ExtraPath == "" {
		if _, err := w.Write([]byte("未设置extrapath目录")); err != nil {
			logger.DefaultLogger.Error(err.Error())
		}
		return
	}

	files := []string{}
	f.updateFile.Range(func(key, value any) bool {
		files = append(files, key.(string))
		return true
	})
	if _, err := w.Write([]byte(strings.Join(files, ","))); err != nil {
		logger.DefaultLogger.Error(err.Error())
	}
}

func (f *FileManager) downloadList(w http.ResponseWriter, r *http.Request) {
	if ServerFlag.ExtraPath == "" {
		if _, err := w.Write([]byte("未设置extrapath目录")); err != nil {
			logger.DefaultLogger.Error(err.Error())
		}
		return
	}

	res, err := fileindex.GetAllFile(ServerFlag.ExtraPath, false)
	if err != nil {
		w.WriteHeader(500)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			logger.DefaultLogger.Error(err.Error())
		}
	} else {
		if _, err := w.Write([]byte(strings.Join(res, ","))); err != nil {
			logger.DefaultLogger.Error(err.Error())
		}
	}
}

// 获取spec信息，在更新中当获取信息时肯定就对应了下载，所以可以在这里将udpateFIle对应的文件删除
func (f *FileManager) downloadSpec(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	if filename == "" {
		if _, err := w.Write([]byte("未设置query参数file")); err != nil {
			logger.DefaultLogger.Error(err.Error())
		}
		return
	}

	targetFile := path.Join(ServerFlag.ExtraPath, filename)
	if !fileindex.IsSubDir(ServerFlag.ExtraPath, targetFile) {
		if _, err := w.Write([]byte("请输入正确的文件名")); err != nil {
			logger.DefaultLogger.Error(err.Error())
		}
		return
	}
	res, err := GetFileSpecInfo(targetFile, int(ServerFlag.ExtraTruncate))
	if err != nil {
		w.WriteHeader(500)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			logger.DefaultLogger.Error(err.Error())
		}
	} else {
		f.updateFile.Delete(filename)
		if _, err := w.Write([]byte(fmt.Sprint(res))); err != nil {
			logger.DefaultLogger.Error(err.Error())
		}
	}
}

func (f *FileManager) downloadMd5(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	md5, err := fileindex.FileMd5BySpec(path.Join(ServerFlag.ExtraPath, filename))
	if err != nil {
		w.WriteHeader(400)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			logger.DefaultLogger.Error(err.Error())
		}
		return
	}

	if _, err := w.Write([]byte(md5)); err != nil {
		logger.DefaultLogger.Error(err.Error())
	}
}

func (f *FileManager) downloadTruncate(w http.ResponseWriter, r *http.Request) {
	filename, trunc := r.URL.Query().Get("file"), r.URL.Query().Get("trunc")
	if filename == "" || trunc == "" {
		if _, err := w.Write([]byte("未设置query参数file或者trunc")); err != nil {
			logger.DefaultLogger.Error(err.Error())
		}
		return
	}

	truncInt, err := strconv.ParseInt(trunc, 10, 64)
	if err != nil {
		w.WriteHeader(500)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			logger.DefaultLogger.Error(err.Error())
		}
		return
	}

	data, err := GetFileData(path.Join(ServerFlag.ExtraPath, filename), truncInt*ServerFlag.ExtraTruncate, int(ServerFlag.ExtraTruncate))
	if err != nil {
		w.WriteHeader(500)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			logger.DefaultLogger.Error(err.Error())
		}
	} else {
		w.Header().Add("Content-Type", "application/offset+octet-stream")
		w.Header().Add("Content-Length", fmt.Sprint(len(data)))
		if _, err := w.Write(data); err != nil {
			logger.DefaultLogger.Error(err.Error())
		}
	}
}

func (f *FileManager) downloadDelete(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	if err := os.Remove(path.Join(ServerFlag.ExtraPath, filename)); err != nil {
		w.WriteHeader(400)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			logger.DefaultLogger.Error(err.Error())
		}
		return
	}

	if _, err := w.Write([]byte("ok")); err != nil {
		logger.DefaultLogger.Error(err.Error())
	}
}
