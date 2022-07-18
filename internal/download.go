package internal

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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
	tempPath  string // chunck保存的临时目录
	size      int64
	offset    int64

	Fingerprint string // download任务的标识
	Metadata    downloadMeta
}

// 先判断当前环境是否已经有这个下载任务了
func NewDownload(fileUrl, fileName, basePath, tempPath string) (*Download, error) {
	if err := os.MkdirAll(basePath, 0o755); err != nil {
		return nil, fmt.Errorf("创建download basepath失败: %w", err)
	}
	if err := os.MkdirAll(path.Join(tempPath, fileName), 0o755); err != nil {
		return nil, fmt.Errorf("创建download temppath失败: %w", err)
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
		tempPath:    tempPath,
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

func (u *Download) ChunckStream(chunck int64) (io.WriteSeeker, error) {
	f, err := os.OpenFile(path.Join(u.tempPath, u.fileName, fmt.Sprintf("%d.chunck", chunck)), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		return nil, fmt.Errorf("创建文件失败: %w", err)
	}

	stream := io.WriteSeeker(f)
	return stream, nil
}

func (u *Download) MergeStream(maxChunck int64) error {
	if _, err := u.stream.Seek(0, io.SeekStart); err != nil {
		return err
	}
	for i := 0; i < int(maxChunck); i++ {
		f, err := os.OpenFile(path.Join(u.tempPath, u.fileName, fmt.Sprintf("%d.chunck", i)), os.O_RDONLY, 0o755)
		if err != nil {
			return fmt.Errorf("打开文件失败: %w", err)
		}
		data, err := ioutil.ReadAll(f)
		if err != nil {
			return fmt.Errorf("读取文件内容失败: %w", err)
		}

		_, err = u.stream.Write(data)
		if err != nil {
			return fmt.Errorf("写入文件失败: %w", err)
		}
	}
	return u.DelTempDir()
}

func (u *Download) DelTempDir() error {
	return os.RemoveAll(path.Join(u.tempPath, u.fileName))
}
