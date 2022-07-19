package internal

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/wwqdrh/fssync/internal/store"
)

type DownloadConfig struct {
	Store   store.DownloadStore
	Resume  bool
	TempDir string
}

type DownloadClient struct {
	Config      *DownloadConfig
	client      *http.Client
	downloadUrl string
}

func NewDownloadClient(url string, config *DownloadConfig) (*DownloadClient, error) {
	return &DownloadClient{
		Config:      config,
		client:      http.DefaultClient,
		downloadUrl: url,
	}, nil
}

// 开始新的下载
func (c *DownloadClient) CreateDownload(download *Download) (*Downloader, error) {
	if err := c.Config.Store.SetOffset(download.Fingerprint, 0); err != nil {
		return nil, err
	}
	return &Downloader{
		client:   c,
		url:      download.fileUrl,
		download: download,
		offset:   0,
		aborted:  false,
	}, nil
}

// 恢复下载
func (c *DownloadClient) ResumeDownload(d *Download) (*Downloader, error) {
	if d == nil {
		return nil, ErrNilUpload
	}

	if !c.Config.Resume {
		return nil, ErrResumeNotEnabled
	} else if len(d.Fingerprint) == 0 {
		return nil, ErrFingerprintNotSet
	}

	offset, found := c.Config.Store.GetOffset(d.Fingerprint)
	if !found {
		return nil, ErrDownloadNotFound
	}

	return &Downloader{
		client:   c,
		url:      d.fileUrl,
		download: d,
		offset:   offset,
		aborted:  false,
	}, nil
}

// 开始新的或者恢复下载
func (c *DownloadClient) CreateOrResumeDownload(d *Download) (*Downloader, error) {
	if d == nil {
		return nil, ErrNilUpload
	}

	uploader, err := c.ResumeDownload(d)
	if err == nil {
		return uploader, err
	} else if (err == ErrResumeNotEnabled) || (err == ErrDownloadNotFound) {
		return c.CreateDownload(d)
	}
	return nil, err
}

func (c *DownloadClient) getmaxChunck(baseurl, filename string) (int64, error) {
	req, err := http.NewRequest("GET", baseurl+"/spec?"+url.Values{"file": []string{filename}}.Encode(), nil)
	if err != nil {
		return -1, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return -1, fmt.Errorf("获取spec失败: %w", err)
	}
	defer res.Body.Close()

	resData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return -1, fmt.Errorf("读取响应失败: %w", err)
	}
	v, err := strconv.ParseInt(string(resData), 10, 64)
	if err != nil {
		return -1, fmt.Errorf("读取响应失败: %w", err)
	}
	return v, nil
}

// 下载切片 fileurl trunc第几个分片
// 保存到临时文件夹下
func (c *DownloadClient) downloadChunck(baseurl, filename string, data io.WriteSeeker, chunck int64) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/truncate?", baseurl)+url.Values{"file": []string{filename}, "trunc": []string{fmt.Sprint(chunck)}}.Encode(), nil)
	if err != nil {
		return err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("下载分片%d失败: %w", chunck, err)
	}
	defer res.Body.Close()

	resData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	_, err = data.Write(resData)
	if err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}
	return nil
}

func (c *DownloadClient) FileList() ([]string, error) {
	req, err := http.NewRequest("GET", c.downloadUrl+"/list", nil)
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(data), ","), nil
}
