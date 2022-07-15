package internal

import (
	"fmt"
	"strings"
)

type downloadMeta map[string]string

type Download struct {
	// stream  io.WriteSeeker
	fileUrl string
	size    int64
	offset  int64

	Fingerprint string // download任务的标识
	Metadata    downloadMeta
}

// 先判断当前环境是否已经有这个下载任务了
func NewDownload(fileUrl string) (*Download, error) {
	metadata := map[string]string{
		"fileurl": fileUrl,
	}
	fingerprint := fileUrl

	return &Download{
		fileUrl:     fileUrl,
		Metadata:    metadata,
		Fingerprint: fingerprint,
	}, nil
}

// Updates the Upload information based on offset.
// func (u *Download) updateProgress(offset int64) {
// 	u.offset = offset
// }

// Returns whether this upload is finished or not.
func (u *Download) Finished() bool {
	return u.offset >= u.size
}

// Returns the progress in a percentage.
func (u *Download) Progress() int64 {
	return (u.offset * 100) / u.size
}

// Returns the current upload offset.
func (u *Download) Offset() int64 {
	return u.offset
}

// Returns the size of the upload body.
func (u *Download) Size() int64 {
	return u.size
}

func (u *Download) EncodedMetadata() string {
	var encoded []string

	for k, v := range u.Metadata {
		encoded = append(encoded, fmt.Sprintf("%s %s", k, b64encode(v)))
	}

	return strings.Join(encoded, ",")
}
