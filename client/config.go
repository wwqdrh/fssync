package client

var ClientDownloadFlag clientDownloadCmdFlag
var ClientUploadFlag clientUploadCmdFlag

type clientDownloadCmdFlag struct {
	Host         string
	DownloadUrl  string
	FileName     string
	DownloadPath string
	SpecPath     string
	TempPath     string // 保存切片的临时目录
	DownAll      bool   // 下载所有的文件
}

type clientUploadCmdFlag struct {
	Host       string
	Uploadfile string
	SpecPath   string
}
