package client

var ClientDownloadFlag clientDownloadCmdFlag
var ClientUploadFlag clientUploadCmdFlag

type clientDownloadCmdFlag struct {
	Host         string
	DownloadUrl  string
	FileName     string
	DownloadPath string
	SpecPath     string
}

type clientUploadCmdFlag struct {
	Host       string
	Uploadfile string
	SpecPath   string
}
