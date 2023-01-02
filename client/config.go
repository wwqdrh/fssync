package client

var ClientDownloadFlag clientDownloadCmdFlag
var ClientUploadFlag clientUploadCmdFlag

type clientDownloadCmdFlag struct {
	DownloadUrl string
	FileName    string

	DownloadPath string
	SpecPath     string
	TempPath     string // 保存切片的临时目录
	DownAll      bool   // 下载所有的文件
	IsDel        bool   // 是否删除
	Watch        bool   // 是否监听更新
	Interval     int    // 更新数据的频率
}

type clientUploadCmdFlag struct {
	Host       string
	Uploadfile string
	SpecPath   string
}

func init() {
	ClientDownloadFlag.DownloadPath = "."
	ClientDownloadFlag.SpecPath = "./tmp/spec"
	ClientDownloadFlag.TempPath = "./tmp/data"
	ClientDownloadFlag.Interval = 10
}
