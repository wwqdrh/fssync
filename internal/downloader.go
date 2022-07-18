package internal

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/wwqdrh/logger"
)

const (
	GoNum  = 10
	MaxErr = 10
)

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
	err = d.client.Config.Store.SetMaxOffset(d.download.Fingerprint, maxTruncate)
	if err != nil {
		return fmt.Errorf("设置最大offset失败: %w", err)
	}

	var wg sync.WaitGroup // 用于等待最后三个
	ch := make(chan struct{}, GoNum)
	errTime := int64(1)
	for errTime < MaxErr {
		val, ok := d.client.Config.Store.GetBlankOffset(d.download.Fingerprint)
		if !ok {
			break
		}

		ch <- struct{}{} // 当写了三次就不能再写了，除非下文有程序执行完了能够继续写入
		wg.Add(1)
		go func(i int64) {
			defer func() {
				wg.Done()
				<-ch
			}()
			logger.DefaultLogger.Info("start chunck: " + fmt.Sprint(i))
			stream, err := d.download.ChunckStream(i)
			if err != nil {
				logger.DefaultLogger.Error("创建stream失败: " + fmt.Sprint(i))
				atomic.AddInt64(&errTime, 1)
				return
			}

			err = d.client.downloadChunck(d.download.fileUrl, d.download.fileName, stream, i)
			if err != nil {
				logger.DefaultLogger.Error("下载chunck失败:" + err.Error())
				if err := d.client.Config.Store.SetFailOffset(d.download.Fingerprint, i); err != nil {
					atomic.AddInt64(&errTime, 1)
					logger.DefaultLogger.Error("写回新offset失败:" + err.Error())
					return
				}
				return
			} else {
				if err := d.client.Config.Store.SetOkOffset(d.download.Fingerprint, i); err != nil {
					atomic.AddInt64(&errTime, 1)
					logger.DefaultLogger.Error("写回新offset失败:" + err.Error())
					return
				}
			}
			logger.DefaultLogger.Info(fmt.Sprintf("chunck: %d downloaded", i))
		}(val)
	}
	wg.Wait()

	if d.client.Config.Store.IsDone(d.download.Fingerprint) {
		fmt.Println("下载完成")
		if err := d.download.MergeStream(maxTruncate); err != nil {
			return err
		}
		fmt.Println("合并完成")
	}

	return nil
}
