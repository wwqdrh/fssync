package server

var ServerFlag serverCmdFlag

type serverCmdFlag struct {
	Port          string
	Store         string
	Urlpath       string
	ExtraPath     string // 额外的直接下载的文件夹
	ExtraTruncate int64  // 额外的直接下载的文件夹 分片的大小
}

func init() {
	ServerFlag.ExtraTruncate = 1 * 1024 * 1024 // 1MB
}
