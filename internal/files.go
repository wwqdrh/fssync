package internal

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// 获取文件的md5字符串
// 根据文件内容进行md5
func FileMd5ByContent(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", fmt.Errorf("FileMd5: 文件打开失败 %w", err)
	}
	defer f.Close()
	md5h := md5.New()
	_, err = io.Copy(md5h, f)
	if err != nil {
		return "", fmt.Errorf("FileMd5: 生成md5失败 %w", err)
	}
	return hex.EncodeToString(md5h.Sum(nil)[:]), nil
}

// 获取文件的md5字符串
// 根据文件名字、文件大小 文件数据的前100byte+后100byte作为md5数据
func FileMd5BySpec(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", fmt.Errorf("FileMd5: 文件打开失败 %w", err)
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return "", fmt.Errorf("FileMd5: 文件stat获取失败 %w", err)
	}

	md5h := md5.New()
	_, err = md5h.Write([]byte(stat.Name()))
	if err != nil {
		return "", fmt.Errorf("FileMd5: 生成md5失败 %w", err)
	}
	_, err = md5h.Write([]byte(fmt.Sprint(stat.Size())))
	if err != nil {
		return "", fmt.Errorf("FileMd5: 生成md5失败 %w", err)
	}
	frontBuf := make([]byte, 100)
	frontNum, err := f.Read(frontBuf)
	if err != nil {
		return "", fmt.Errorf("FileMd5: 生成md5失败 %w", err)
	}
	_, err = md5h.Write(frontBuf[:frontNum])
	if err != nil {
		return "", fmt.Errorf("FileMd5: 生成md5失败 %w", err)
	}
	backBuf := make([]byte, 100)
	_, err = f.Seek(100, os.SEEK_END)
	if err != nil {
		return "", fmt.Errorf("FileMd5: 生成md5失败 %w", err)
	}
	backNum, err := f.Read(backBuf)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("FileMd5: 生成md5失败 %w", err)
	}
	_, err = md5h.Write(backBuf[:backNum])
	if err != nil {
		return "", fmt.Errorf("FileMd5: 生成md5失败 %w", err)
	}

	return hex.EncodeToString(md5h.Sum(nil)), nil
}
