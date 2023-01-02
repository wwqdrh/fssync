package server

var ServerFlag serverCmdFlag

type serverCmdFlag struct {
	Port          int
	Store         string
	Urlpath       string
	ExtraPath     string // 额外的直接下载的文件夹
	ExtraTruncate int64  // 额外的直接下载的文件夹 分片的大小
}

func init() {
	ServerFlag.Port = 1080
	ServerFlag.ExtraTruncate = 1 * 1024 * 1024 // 1MB
	ServerFlag.Urlpath = "/files"
	ServerFlag.ExtraPath = "."
	ServerFlag.Store = "./upload"
}
