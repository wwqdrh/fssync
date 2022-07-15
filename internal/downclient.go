package internal

import "github.com/wwqdrh/fssync/internal/store"

type DownloadConfig struct {
	Store  store.DownloadStore
	Resume bool
}

type DownloadClient struct {
	Config *DownloadConfig
}

// 开始新的下载
func (c *DownloadClient) CreateDownload(download *Download) (*Downloader, error) {
	if err := c.Config.Store.SetOffset(download.Fingerprint, 0); err != nil {
		return nil, err
	}
	return &Downloader{}, nil
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

	return &Downloader{}, nil
}

// 开始新的或者恢复下载
func (c *DownloadClient) CreateOrResumeDownload(d *Download) (*Downloader, error) {
	if d == nil {
		return nil, ErrNilUpload
	}

	uploader, err := c.ResumeDownload(d)
	if err == nil {
		return uploader, err
	} else if (err == ErrResumeNotEnabled) || (err == ErrUploadNotFound) {
		return c.CreateDownload(d)
	}
	return nil, err
}
