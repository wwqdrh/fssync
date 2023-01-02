package store

type UploadStore interface {
	Get(fingerprint string) (string, bool)
	Set(fingerprint, url string) error
	Delete(fingerprint string) error
	Close() error
}

// 需要支持多协程，因此chunk需要支持非连续机制
type DownloadStore interface {
	GetOffset(figerprint string) (int64, bool) // 获取当前最小的offset
	SetOffset(figerprint string, offset int64) error
	Delete(fingerprint string) error
	Close() error

	SetMaxOffset(figerprint string, offset int64) error // 设置一个figerprint最大的切片数
	GetMaxOffset(figerprint string) (int64, bool)
	GetBlankOffset(figerprint string) (int64, bool)    // 并发安全 获取一个还未下载的切片
	SetOkOffset(figerprint string, offset int64) error // 并发安全 标记一个切片已经下载完成
	SetFailOffset(figerprint string, offset int64) error
	IsDone(figerprint string) bool // 判断是否下载完成
	IsCombile(figerprint string) error
	SetCombile(figerprint string) error
}
