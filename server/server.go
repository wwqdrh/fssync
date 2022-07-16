package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strconv"
	"strings"
	"syscall"

	"github.com/tus/tusd/pkg/filestore"
	tusd "github.com/tus/tusd/pkg/handler"
	"github.com/wwqdrh/logger"
)

func Start(ctx context.Context) error {
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
				logger.DefaultLogger.Info(fmt.Sprintf("Upload %s finished\n", event.Upload.ID))
			case <-ctx.Done():
				return
			}
		}
	}()

	http.Handle(ServerFlag.Urlpath, http.StripPrefix(ServerFlag.Urlpath, handler))
	registerAPI()

	go func() {
		logger.DefaultLogger.Info(ServerFlag.Port)
		err = http.ListenAndServe(ServerFlag.Port, nil)
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

////////////////////
// extra api
////////////////////
func registerAPI() {
	http.HandleFunc("/download/list", downloadList)
	http.HandleFunc("/download/spec", downloadSpec)
	http.HandleFunc("/download/truncate", downloadTruncate)
}

func downloadList(w http.ResponseWriter, r *http.Request) {
	if ServerFlag.ExtraPath == "" {
		w.Write([]byte("未设置extrapath目录"))
		return
	}

	res, err := ListDirFile(ServerFlag.ExtraPath, false)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	} else {
		w.Write([]byte(strings.Join(res, ",")))
	}
}

func downloadSpec(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	if filename == "" {
		w.Write([]byte("未设置query参数file"))
		return
	}

	res, err := GetFileSpecInfo(path.Join(ServerFlag.ExtraPath, filename), int(ServerFlag.ExtraTruncate))
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	} else {
		w.Write([]byte(fmt.Sprint(res)))
	}
}

func downloadTruncate(w http.ResponseWriter, r *http.Request) {
	filename, trunc := r.URL.Query().Get("file"), r.URL.Query().Get("trunc")
	if filename == "" || trunc == "" {
		w.Write([]byte("未设置query参数file或者trunc"))
		return
	}

	truncInt, err := strconv.ParseInt(trunc, 10, 64)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	data, err := GetFileData(path.Join(ServerFlag.ExtraPath, filename), truncInt*ServerFlag.ExtraTruncate, int(ServerFlag.ExtraTruncate))
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	} else {
		w.Write(data)
	}
}
