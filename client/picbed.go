package client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"golang.design/x/clipboard"
)

func fnUpload(prefix, file, cookie string) (string, string, error) {
	records, err := getSpec()
	if err != nil {
		return "", "", err
	}
	if cookie == "" {
		cookie = records.Cookie
	}
	if prefix == "" {
		prefix = records.Prefix
	}
	uploader := NewWeiboUploader()
	fileURL := file
	domain := prefix + "https://cdn.ipfsscan.io/weibo"
	quality := "large"
	cookieMode := true

	return uploader.Upload(fileURL, "", "", cookie, domain, quality, cookieMode)
}

type picSpec struct {
	Cookie string                `json:"cookie"`
	Prefix string                `json:"prefix"`
	Data   map[string]*picRecord `json:"data"`
}

type picRecord struct {
	Url  string `json:"url"`
	Path string `json:"path"`
}

func getSpec() (*picSpec, error) {
	specPath := path.Join(RootSpecPath, "picbed.json")
	var records picSpec
	if data, err := os.ReadFile(specPath); err == nil {
		if err := json.Unmarshal(data, &records); err != nil {
			return nil, err
		}
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	return &records, nil
}

func addRecord(id, local, url string) error {
	records, err := getSpec()
	if err != nil {
		return err
	}

	records.Data[id] = &picRecord{
		Url:  url,
		Path: local,
	}
	content, err := json.Marshal(records)
	if err != nil {
		return err
	}

	return os.WriteFile(path.Join(RootSpecPath, "picbed.json"), content, os.ModePerm)
}

// WeiboUploader 用于上传文件到微博
type WeiboUploader struct {
	URL            string
	FileExtensions []string
}

// NewWeiboUploader 创建一个新的 WeiboUploader 实例
func NewWeiboUploader() *WeiboUploader {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	return &WeiboUploader{
		URL:            "https://picupload.weibo.com/interface/pic_upload.php?ori=1&mime=image%2Fjpeg&data=base64&url=0&markpos=1&logo=&nick=0&marks=1&app=miniblog",
		FileExtensions: []string{"jpeg", "jpg", "png", "gif", "bmp"},
	}
}

// Upload 上传文件到微博
// fileURL, 如果是空字符串，那么就尝试从剪切版中获取数据
func (u *WeiboUploader) Upload(fileURL string, username, password, cookie, domain, quality string, cookieMode bool) (string, string, error) {
	if cookieMode && cookie == "" {
		return "", "", fmt.Errorf("there is a problem with the map bed configuration, please check")
	} else if !cookieMode && (username == "" || password == "") {
		return "", "", fmt.Errorf("there is a problem with the map bed configuration, please check")
	}

	var loginCookie string
	if !cookieMode {
		cookie, err := u.login(username, password)
		if err != nil {
			return "", "", err
		}
		loginCookie = cookie
	} else {
		loginCookie = cookie
	}

	var fileData []byte
	var err error
	var fileExtension string = ".jpg"

	if fileURL != "" {
		fileData, err = os.ReadFile(fileURL)
		if err != nil {
			return "", "", err
		}

		// fileName := filepath.Base(fileURL)
		fileExtension = strings.ToLower(filepath.Ext(fileURL))
		if fileExtension == ".gif" {
			fileExtension = ".gif"
		} else {
			fileExtension = ".jpg"
		}
	} else {
		fileData = clipboard.Read(clipboard.FmtImage)
	}

	fileBase64 := base64.StdEncoding.EncodeToString(fileData)

	localPath, err := u.writeLocal(fileData, fileBase64)
	if err != nil {
		return "", "", err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("b64_data", fileBase64)
	writer.Close()

	req, err := http.NewRequest("POST", u.URL, body)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Cookie", loginCookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	pidPid, err := parsePicPid(string(bodyBytes))
	if err != nil {
		return "", "", err
	}

	url := fmt.Sprintf("%s/%s/%s%s", domain, quality, pidPid, fileExtension)
	return url, localPath, nil
}

func (u *WeiboUploader) writeLocal(data []byte, base64 string) (string, error) {
	length := len(base64)
	dirPath := filepath.Join(RootSpecPath, fmt.Sprint(length), string([]rune(base64)[length-20:length-10]))

	// 创建目录(包括任何必需的父目录)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	// 构建完整的文件路径
	filePath := filepath.Join(dirPath, string([]rune(base64)[length-10:])+".jpg")

	// 将图片数据写入文件
	return filePath, os.WriteFile(filePath, data, 0644)
}

func (u *WeiboUploader) login(username, password string) (string, error) {
	loginURL := "https://passport.weibo.cn/sso/login"

	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)

	req, err := http.NewRequest("POST", loginURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", loginURL)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var respData map[string]interface{}
	err = json.Unmarshal(bodyBytes, &respData)
	if err != nil {
		return "", err
	}

	retcode, ok := respData["retcode"].(float64)
	if !ok || int(retcode) != 20000000 {
		errMsg, _ := respData["msg"].(string)
		return "", fmt.Errorf(errMsg)
	}

	cookies := resp.Header.Get("Set-Cookie")
	return cookies, nil
}

func parsePicPid(respString string) (string, error) {
	start := strings.Index(respString, `"pid":"`) + 7
	end := strings.Index(respString[start:], `"`) + start
	if end == -1 {
		return "", fmt.Errorf("unable to parse pic pid from response")
	}
	return respString[start:end], nil
}
