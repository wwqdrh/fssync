package server

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/wwqdrh/logger"
)

// source: ./testdata or testdata
func ListDirFile(source string, prefix bool) ([]string, error) {
	source = strings.TrimLeft(source, "./")
	dirStack := []string{source}

	res := []string{}
	for len(dirStack) > 0 {
		cur := dirStack[0]
		dirStack = dirStack[1:]

		files, err := ioutil.ReadDir(cur)
		if err != nil {
			logger.DefaultLogger.Warn(cur + " 不是文件夹")
			continue
		}

		for _, item := range files {
			if item.IsDir() {
				dirStack = append(dirStack, path.Join(cur, item.Name()))
			} else {
				res = append(res, path.Join(cur, item.Name()))
			}
		}
	}

	if !prefix {
		for i := 0; i < len(res); i++ {
			cur := strings.TrimPrefix(res[i], source)
			cur = strings.TrimPrefix(cur, "/")
			res[i] = cur
		}
	}

	return res, nil
}

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
