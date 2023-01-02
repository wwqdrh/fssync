package server

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/wwqdrh/gokit/logger"
)

// 目标文件 分片大小
func GetFileSpecInfo(source string, truncate int) (int64, error) {
	f, err := os.Open(source)
	if err != nil {
		return -1, fmt.Errorf("打开文件失败: %w", err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return -1, fmt.Errorf("获取文件信息失败: %w", err)
	}

	return (fi.Size() + int64(truncate) - 1) / int64(truncate), nil
}

// 获取目标文件指定部分的数据
func GetFileData(source string, offset int64, trucate int) ([]byte, error) {
	f, err := os.Open(source)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer f.Close()

	reader := func(r io.Reader) io.Reader {
		return r
	}(f)

	stream, ok := reader.(io.ReadSeeker)
	if !ok {
		logger.DefaultLogger.Warn(source + "文件没有实现io.ReadSeeker接口")
		buf := new(bytes.Buffer)
		_, err := buf.ReadFrom(reader)
		if err != nil {
			return nil, err
		}
		stream = bytes.NewReader(buf.Bytes())
	}

	data := make([]byte, trucate)
	_, err = stream.Seek(offset, 0)
	if err != nil {
		return nil, fmt.Errorf("移动seek失败: %w", err)
	}
	n, err := stream.Read(data)
	if err != nil {
		return nil, fmt.Errorf("读取数据失败: %w", err)
	}
	return data[:n], nil
}
