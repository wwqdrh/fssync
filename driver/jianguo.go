package driver

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type JianguoDriver struct {
	entry        string
	authName     string
	authPassword string
}

// 定义结构体与XML节点对应
type Response struct {
	Href     string `xml:"href"`
	Propstat Propstat
}

type Propstat struct {
	Prop   Prop   `xml:"prop"`
	Status string `xml:"status"`
}

type Prop struct {
	GetContentType          string                  `xml:"getcontenttype"`
	DisplayName             string                  `xml:"displayname"`
	Owner                   string                  `xml:"owner"`
	ResourceType            ResourceType            `xml:"resourcetype"`
	GetContentLength        int64                   `xml:"getcontentlength"`
	GetLastModified         string                  `xml:"getlastmodified"`
	CurrentUserPrivilegeSet CurrentUserPrivilegeSet `xml:"current-user-privilege-set"`
}

type ResourceType struct {
	Collection interface{} `xml:"collection"`
}

type CurrentUserPrivilegeSet struct {
	Privileges []Privilege `xml:"privilege"`
}

type Privilege struct {
	Name string `xml:",chardata"`
}

type Multistatus struct {
	XMLName   xml.Name   `xml:"multistatus"`
	Responses []Response `xml:"response"`
}

func NewJianguoDriver() IDriver {
	return &JianguoDriver{
		entry: "https://dav.jianguoyun.com/dav/我的坚果云/",
	}
}

func (d *JianguoDriver) Auth(name, password string) {
	d.authName = name
	d.authPassword = password
}

func (d *JianguoDriver) IsAuth() bool {
	return d.authName != "" && d.authPassword != ""
}

func (d *JianguoDriver) Download(url string) error {
	if !d.IsAuth() {
		return ErrNotAuth
	}
	if url == "" {
		url = d.entry
	} else {
		url = d.entry + url
	}
	// 创建HTTP请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// 添加身份验证凭据(如果需要)
	req.SetBasicAuth(d.authName, d.authPassword)

	// 发送HTTP请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error downloading file: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error downloading file. Status code: %d", resp.StatusCode)
	}

	// 创建本地文件
	outFile, err := os.Create("downloaded_file.txt")
	if err != nil {
		return fmt.Errorf("error creating local file: %w", err)
	}
	defer outFile.Close()

	// 将响应体写入本地文件
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("error writing to local file: %w", err)
	}

	return nil
}

func (d *JianguoDriver) List(url string) ([]FileItem, error) {
	if !d.IsAuth() {
		return nil, ErrNotAuth
	}
	if url == "" {
		url = d.entry
	} else {
		url = d.entry + url
	}
	// 发送PROPFIND请求以获取文件列表
	req, _ := http.NewRequest("PROPFIND", url, nil)
	req.Header.Set("Depth", "1")
	req.SetBasicAuth(d.authName, d.authPassword)

	// 发送请求并获取响应
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	// 解析XML响应
	ms := Multistatus{}
	err = xml.Unmarshal(body, &ms)
	if err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}

	res := []FileItem{}
	for _, resp := range ms.Responses {
		item := FileItem{
			Name:          resp.Propstat.Prop.DisplayName,
			Href:          resp.Href,
			Owner:         resp.Propstat.Prop.Owner,
			Status:        resp.Propstat.Status,
			ResourceType:  resp.Propstat.Prop.ResourceType.Collection,
			ContentType:   resp.Propstat.Prop.GetContentType,
			ContentLength: resp.Propstat.Prop.GetContentLength,
			LastModify:    resp.Propstat.Prop.GetLastModified,
		}
		privilege := []string{}
		for _, priv := range resp.Propstat.Prop.CurrentUserPrivilegeSet.Privileges {
			privilege = append(privilege, priv.Name)
		}
		item.Privileges = privilege
		res = append(res, item)
	}
	return res, nil
}

func (d *JianguoDriver) Delete(remote string) error {
	if !d.IsAuth() {
		return ErrNotAuth
	}
	if remote == "" {
		return ErrInvUrl
	}

	req, err := http.NewRequest("DELETE", d.entry+remote, nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(d.authName, d.authPassword)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete file: %s", resp.Status)
	}

	return nil
}

func (d *JianguoDriver) Update(local, remote string) error {
	if !d.IsAuth() {
		return ErrNotAuth
	}
	if remote == "" {
		return ErrInvUrl
	}

	d.createDirectories(local)

	file, err := os.Open(local)
	if err != nil {
		return err
	}
	defer file.Close()

	req, err := http.NewRequest("PUT", d.entry+remote, file)
	if err != nil {
		return err
	}

	req.SetBasicAuth(d.authName, d.authPassword)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to upload file: %s", resp.Status)
	}
	return nil
}

func (d *JianguoDriver) createDirectories(local string) error {
	if !d.IsAuth() {
		return ErrNotAuth
	}
	dirPath := filepath.Dir(local)
	if dirPath == "." || dirPath == "/" || dirPath == "" {
		return nil
	}

	pathParts := strings.Split(dirPath, "/")
	currentPath := ""
	for _, part := range pathParts {
		if part == "" {
			continue
		}

		currentPath = filepath.Join(currentPath, part)

		req, err := http.NewRequest("MKCOL", d.entry+currentPath, nil)
		if err != nil {
			return err
		}
		req.SetBasicAuth(d.authName, d.authPassword)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
			return fmt.Errorf("failed to upload file: %s", resp.Status)
		}
	}

	return nil
}
