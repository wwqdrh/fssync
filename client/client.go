package client

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/wwqdrh/fssync/client/download"
	"github.com/wwqdrh/fssync/client/upload"
	"github.com/wwqdrh/fssync/pkg/store"
	"github.com/wwqdrh/gokit/logger"
)

func newClient() (*download.DownloadClient, error) {
	if err := os.MkdirAll(ClientDownloadFlag.SpecPath, 0o777); err != nil {
		return nil, fmt.Errorf("创建spec失败: %w", err)
	}
	if err := os.MkdirAll(ClientDownloadFlag.TempPath, 0o777); err != nil {
		return nil, fmt.Errorf("创建temp失败: %w", err)
	}

	s, err := store.NewLeveldbStore(ClientDownloadFlag.SpecPath)
	if err != nil {
		return nil, fmt.Errorf("持久化组件初始化失败: %w", err)
	}
	v, ok := s.(store.DownloadStore)
	if !ok {
		return nil, fmt.Errorf("持久化组件初始化失败: %w", errors.New("leveldb store未实现uploadStore接口"))
	}
	// defer v.Close()

	client, err := download.NewDownloadClient(ClientDownloadFlag.DownloadUrl, &download.DownloadConfig{
		Resume:  true,
		Store:   v,
		TempDir: ClientDownloadFlag.TempPath,
	})
	if err != nil {
		return nil, fmt.Errorf("tus client初始化失败: %w", err)
	}
	return client, nil
}

func UploadStart() error {
	f, err := os.Open(ClientUploadFlag.Uploadfile)
	if err != nil {
		return fmt.Errorf("打开目标文件失败: %w", err)
	}
	defer f.Close()

	s, err := store.NewLeveldbStore(ClientUploadFlag.SpecPath)
	if err != nil {
		return fmt.Errorf("持久化组件初始化失败: %w", err)
	}
	v, ok := s.(store.UploadStore)
	if !ok {
		return fmt.Errorf("持久化组件初始化失败: %w", errors.New("leveldb store未实现uploadStore接口"))
	}
	defer v.Close()

	client, err := upload.NewUploadClient(ClientUploadFlag.Host, &upload.UploadConfig{
		ChunkSize:           2 * 1024 * 1024,
		Resume:              true,
		OverridePatchMethod: true,
		Store:               v,
		Header:              make(http.Header),
		HttpClient:          nil,
	})
	if err != nil {
		return fmt.Errorf("tus client初始化失败: %w", err)
	}
	up, err := upload.NewUploadFromFile(f)
	if err != nil {
		return fmt.Errorf("tus client初始化文件上传失败: %w", err)
	}

	uploader, err := client.CreateOrResumeUpload(up)
	if err != nil {
		return fmt.Errorf("tus client初始化文件上传失败: %w", err)
	}
	err = uploader.Upload()
	if err != nil {
		return fmt.Errorf("tus client文件上传失败: %w", err)
	}

	return nil
}

func DownloadStart() error {
	ClientDownloadFlag.Init()

	client, err := newClient()
	if err != nil {
		return err
	}
	defer client.Close()

	if ClientDownloadFlag.DownAll {
		fileList, err := client.FileList()
		if err != nil {
			return fmt.Errorf("获取下载列表失败: %w", err)
		}
		if err := DownloadAll(client, fileList, false); err != nil {
			return err
		}

		if ClientDownloadFlag.Watch {
			timer := time.NewTicker(time.Duration(ClientDownloadFlag.Interval) * time.Second)
			defer timer.Stop()
			for range timer.C {
				fileList, err = client.FileUpdateList()
				if err != nil {
					logger.DefaultLogger.Warn(err.Error())
					continue
				}
				logger.DefaultLogger.Infox("update file: %s", nil, strings.Join(fileList, ","))
				if len(fileList) == 0 {
					continue
				}

				if err := DownloadAll(client, fileList, true); err != nil {
					logger.DefaultLogger.Warn(err.Error())
				}
			}
		}

		return nil
	} else {
		return downloadOne(client, ClientDownloadFlag.FileName, false)
	}
}

func DownloadAll(client *download.DownloadClient, fileList []string, force bool) error {
	for _, item := range fileList {
		if err := downloadOne(client, strings.TrimSpace(item), force); err != nil {
			logger.DefaultLogger.Warn(err.Error())
		} else {
			logger.DefaultLogger.Debug(item + "下载成功")
		}
	}
	return nil
}

func DownloadStartByList(files ...string) []error {
	client, err := newClient()
	if err != nil {
		return []error{err}
	}
	defer client.Close()

	errs := []error{}
	for _, item := range files {
		if err := downloadOne(client, item, false); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func downloadOne(client *download.DownloadClient, fileName string, force bool) error {
	if fileName == "" {
		return errors.New("文件名不能为空")
	}

	if err := os.MkdirAll(path.Dir(path.Join(ClientDownloadFlag.DownloadPath, fileName)), os.ModePerm); err != nil {
		return err
	}

	download, err := download.NewDownload(ClientDownloadFlag.DownloadUrl, fileName, ClientDownloadFlag.DownloadPath, ClientDownloadFlag.TempPath)
	if err != nil {
		return fmt.Errorf("newdownload err %w", err)
	}

	downloader, err := client.CreateOrResumeDownload(download, force)
	if err != nil {
		return fmt.Errorf("downlaoded err: %w", err)
	}
	err = downloader.Download(ClientDownloadFlag.IsDel)
	if err != nil {
		return fmt.Errorf("downloaded merge err: %w", err)
	}
	return nil
}

func DownloadList() ([]string, error) {
	if err := os.MkdirAll(ClientDownloadFlag.SpecPath, 0o777); err != nil {
		return nil, fmt.Errorf("创建spec失败: %w", err)
	}
	if err := os.MkdirAll(ClientDownloadFlag.TempPath, 0o777); err != nil {
		return nil, fmt.Errorf("创建temp失败: %w", err)
	}

	s, err := store.NewLeveldbStore(ClientDownloadFlag.SpecPath)
	if err != nil {
		return nil, fmt.Errorf("持久化组件初始化失败: %w", err)
	}
	v, ok := s.(store.DownloadStore)
	if !ok {
		return nil, fmt.Errorf("持久化组件初始化失败: %w", errors.New("leveldb store未实现uploadStore接口"))
	}
	defer v.Close()

	client, err := download.NewDownloadClient(ClientDownloadFlag.DownloadUrl, &download.DownloadConfig{
		Resume:  true,
		Store:   v,
		TempDir: ClientDownloadFlag.TempPath,
	})
	if err != nil {
		return nil, fmt.Errorf("tus client初始化失败: %w", err)
	}

	return client.FileList()
}
