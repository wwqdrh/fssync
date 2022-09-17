package download

import (
	"errors"
	"fmt"
)

var (
	ErrNilUpload         = errors.New("upload can't be nil")
	ErrLargeUpload       = errors.New("upload body is to large")
	ErrUploadNotFound    = errors.New("upload not found")
	ErrResumeNotEnabled  = errors.New("resuming not enabled")
	ErrFingerprintNotSet = errors.New("fingerprint not set")
	ErrNilDownload       = errors.New("download can't be nil")
	ErrDownloadNotFound  = errors.New("download not found")
)

type ClientError struct {
	Code int
	Body []byte
}

func (c ClientError) Error() string {
	return fmt.Sprintf("unexpected status code: %d", c.Code)
}
