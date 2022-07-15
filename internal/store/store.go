package store

type UploadStore interface {
	Get(fingerprint string) (string, bool)
	Set(fingerprint, url string) error
	Delete(fingerprint string) error
	Close() error
}

type DownloadStore interface {
	GetOffset(figerprint string) (int64, bool)
	SetOffset(figerprint string, offset int64) error
	Delete(fingerprint string)
	Close()
}
