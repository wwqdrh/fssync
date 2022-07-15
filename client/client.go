package client

import (
	"fmt"
	"net/http"
	"os"

	"github.com/wwqdrh/fssync/internal"
	"github.com/wwqdrh/fssync/internal/store"
)

func Start() error {
	f, err := os.Open(ClientFlag.Uploadfile)
	if err != nil {
		return fmt.Errorf("打开目标文件失败: %w", err)
	}
	defer f.Close()

	s, err := store.NewLeveldbStore(ClientFlag.SpecPath)
	if err != nil {
		return fmt.Errorf("持久化组件初始化失败: %w", err)
	}

	client, err := internal.NewClient(ClientFlag.Host, &internal.Config{
		ChunkSize:           2 * 1024 * 1024,
		Resume:              true,
		OverridePatchMethod: true,
		Store:               s,
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

	uploader, err := client.CreateUpload(upload)
	if err != nil {
		return fmt.Errorf("tus client初始化文件上传失败: %w", err)
	}
	err = uploader.Upload()
	if err != nil {
		return fmt.Errorf("tus client文件上传失败: %w", err)
	}

	return nil
}
