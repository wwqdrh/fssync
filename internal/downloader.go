package internal

import "fmt"

type Downloader struct {
	client   *DownloadClient
	url      string
	download *Download
	offset   int64
	aborted  bool
}

func (d *Downloader) Abort() {
	d.aborted = true
}

func (d *Downloader) IsAborted() bool {
	return d.aborted
}

func (d *Downloader) Offset() int64 {
	return d.offset
}

// 查看
func (d *Downloader) Download() error {
	maxTruncate, err := d.client.getmaxChunck(d.download.fileUrl, d.download.fileName)
	if err != nil {
		return err
	}
	for d.offset < maxTruncate && !d.aborted {
		err := d.client.downloadChunck(d.download.fileUrl, d.download.fileName, d.download.stream, d.offset)
		if err != nil {
			return err
		}
		d.offset += 1
		if err := d.client.Config.Store.SetOffset(d.download.Fingerprint, d.offset); err != nil {
			return fmt.Errorf("写回新offset失败: %w", err)
		}
	}
	return nil
}
