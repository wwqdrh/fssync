package client

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/wwqdrh/fssync/internal"
	"github.com/wwqdrh/fssync/internal/store"
)

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

	client, err := internal.NewUploadClient(ClientUploadFlag.Host, &internal.UploadConfig{
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
	upload, err := internal.NewUploadFromFile(f)
	if err != nil {
		return fmt.Errorf("tus client初始化文件上传失败: %w", err)
	}

	uploader, err := client.CreateOrResumeUpload(upload)
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
	if err := os.MkdirAll(ClientDownloadFlag.SpecPath, 0o755); err != nil {
		return fmt.Errorf("创建spec失败: %w", err)
	}
	s, err := store.NewLeveldbStore(ClientDownloadFlag.SpecPath)
	if err != nil {
		return fmt.Errorf("持久化组件初始化失败: %w", err)
	}
	v, ok := s.(store.DownloadStore)
	if !ok {
		return fmt.Errorf("持久化组件初始化失败: %w", errors.New("leveldb store未实现uploadStore接口"))
	}
	defer v.Close()

	client, err := internal.NewDownloadClient(ClientDownloadFlag.DownloadUrl, &internal.DownloadConfig{
		Resume: true,
		Store:  v,
	})
	if err != nil {
		return fmt.Errorf("tus client初始化失败: %w", err)
	}

	download, err := internal.NewDownload(ClientDownloadFlag.DownloadUrl, ClientDownloadFlag.FileName, ClientDownloadFlag.DownloadPath)
	if err != nil {
		return fmt.Errorf("tus client初始化文件上传失败: %w", err)
	}

	downloader, err := client.CreateOrResumeDownload(download)
	if err != nil {
		return fmt.Errorf("tus client初始化文件上传失败: %w", err)
	}
	err = downloader.Download()
	if err != nil {
		return fmt.Errorf("tus client文件上传失败: %w", err)
	}
	return nil
}
