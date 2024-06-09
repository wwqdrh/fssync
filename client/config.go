package client

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/wwqdrh/clipboard"
	"github.com/wwqdrh/gokit/clitool"
	"github.com/wwqdrh/gokit/logger"
)

var RootSpecPath = ""

var (
	ClientDownloadFlag = ClientDownloadCmdFlag{
		DownloadUrl:  "http://127.0.0.1:1080",
		DownloadPath: ".",
		SpecPath:     "./tmp/spec",
		TempPath:     "./tmp/data",
		Interval:     10,
	}
	ClientUploadFlag ClientUploadCmdFlag
	ClientWebDavFlag = struct {
		CfgFile  string   `name:"cfg" desc:"存储配置文件地址"`
		Name     string   `name:"name" desc:"webdav服务名, example(坚果云)"`
		Ignores  []string `name:"ignores" desc:"忽略部分文件不上传"`
		Work     string   `name:"work" desc:"工作目录"`
		Interval int      `name:"interval" desc:"更新检查频率(ms)"`
	}{
		CfgFile:  "config.json",
		Name:     "坚果云",
		Work:     ".",
		Interval: 10000,
		Ignores:  []string{"config.json"},
	}

	ClientPicBedFlag = struct {
		Prefix string `name:"prefix"`
		PicId  string `name:"id"`
		File   string `name:"file"`
		Cookie string `name:"cookie"`
	}{}
)

type ClientDownloadCmdFlag struct {
	DownloadUrl  string `name:"host" desc:"服务端地址"`
	FileName     string `name:"file" desc:"需要下载的单个文件的名字"`
	DownloadPath string `name:"down" desc:"下载文件保存路径"`
	SpecPath     string `name:"spec" desc:"切片信息等保存路径"`
	TempPath     string `name:"temp" desc:"保存切片的临时目录"`             // 保存切片的临时目录
	DownAll      bool   `name:"all" desc:"是否下载所有的文件"`              // 下载所有的文件
	IsDel        bool   `name:"del" desc:"是否删除"`                   // 是否删除
	Watch        bool   `name:"watch" desc:"是监听更新"`                // 是否监听更新
	Interval     int    `name:"interval" alias:"i" desc:"更新数据的频率"` // 更新数据的频率
}

type ClientUploadCmdFlag struct {
	Host       string
	Uploadfile string
	SpecPath   string
}

func (c *ClientDownloadCmdFlag) Init() {
	var err error

	if c.DownloadPath, err = filepath.Abs(c.DownloadPath); err != nil {
		logger.DefaultLogger.Fatal(err.Error())
	}
	if c.SpecPath, err = filepath.Abs(c.SpecPath); err != nil {
		logger.DefaultLogger.Fatal(err.Error())
	}
	if c.TempPath, err = filepath.Abs(c.TempPath); err != nil {
		logger.DefaultLogger.Fatal(err.Error())
	}
}

func Command() *clitool.Command {
	cmd := &clitool.Command{
		Cmd: &cobra.Command{
			Use:   "client",
			Short: "client",
		},
	}
	cmd.Add(&clitool.Command{
		Cmd: &cobra.Command{
			Use:   "webdav",
			Short: "webdav",
			RunE: func(cmd *cobra.Command, args []string) error {
				return SyncWebdav()
			},
		},
		Values: &ClientWebDavFlag,
	})
	cmd.Add(&clitool.Command{
		Cmd: &cobra.Command{
			Use: "pic",
			RunE: func(cmd *cobra.Command, args []string) error {
				picid := ClientPicBedFlag.PicId
				if picid == "" {
					picid = time.Now().Format("20060102150405")
				}
				picurl, localPath, err := fnUpload(ClientPicBedFlag.Prefix, ClientPicBedFlag.File, ClientPicBedFlag.Cookie)
				if err != nil {
					return err
				}
				if err := addRecord(picid, localPath, picurl); err != nil {
					return err
				}
				fmt.Printf("data-id=\"%s\"\n", picid)
				fmt.Printf("url=\"%s\"\n", picurl)
				return clipboard.WriteAll(picurl)
			},
		},
		Values: &ClientPicBedFlag,
	})
	cmd.Add(&clitool.Command{
		Cmd: &cobra.Command{
			Use:   "download",
			Short: "download",
			RunE: func(cmd *cobra.Command, args []string) error {
				return DownloadStart()
			},
		},
		Values: &ClientDownloadFlag,
	})
	return cmd
}
