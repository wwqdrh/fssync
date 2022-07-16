package internal

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

type downloadMeta map[string]string

type Download struct {
	stream    io.WriteSeeker
	fileUrl   string
	fileName  string
	localPath string // 下载文件存储的位置
	size      int64
	offset    int64

	Fingerprint string // download任务的标识
	Metadata    downloadMeta
}

// 先判断当前环境是否已经有这个下载任务了
func NewDownload(fileUrl, fileName, basePath string) (*Download, error) {
	if err := os.MkdirAll(basePath, 0o777); err != nil {
		return nil, fmt.Errorf("创建download basepath失败: %w", err)
	}

	metadata := map[string]string{
		"fileurl":  fileUrl,
		"filename": fileName,
		"basepath": basePath,
	}
	fingerprint := fmt.Sprintf("%s-%s-%s", fileUrl, fileName, basePath)

	stream, err := newWriterStream(path.Join(basePath, fileName))
	if err != nil {
		return nil, fmt.Errorf("创建download stream失败: %w", err)
	}
	return &Download{
		stream:      stream,
		fileUrl:     fileUrl,
		fileName:    fileName,
		localPath:   path.Join(basePath, fileName),
		Metadata:    metadata,
		Fingerprint: fingerprint,
	}, nil
}

func newWriterStream(source string) (io.WriteSeeker, error) {
	f, err := os.OpenFile(source, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		return nil, fmt.Errorf("创建文件失败: %w", err)
	}

	fw := func(w io.Writer) io.Writer { return w }(f)
	stream, ok := fw.(io.WriteSeeker)
	if !ok {
		return nil, errors.New("wrietseeker构造失败")
	}
	return stream, nil
}

func (u *Download) Finished() bool {
	return u.offset >= u.size
}

func (u *Download) Progress() int64 {
	return (u.offset * 100) / u.size
}

func (u *Download) Offset() int64 {
	return u.offset
}

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
