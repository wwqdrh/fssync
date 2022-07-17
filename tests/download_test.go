package tests

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/wwqdrh/fssync/client"
	"github.com/wwqdrh/fssync/server"
)

type DownloadSuite struct {
	suite.Suite

	cancel context.CancelFunc
}

func init() {
	server.ServerFlag.Port = ":1080"
	server.ServerFlag.Store = "./testdata/store"
	server.ServerFlag.Urlpath = "/files/"
	server.ServerFlag.ExtraPath = "./testdata/downloadextra"

	client.ClientUploadFlag.Host = "http://127.0.0.1:1080/files/"
	client.ClientUploadFlag.SpecPath = "./testdata/uploadspec"
}

func TestDownloadSuite(t *testing.T) {
	suite.Run(t, &DownloadSuite{})
}

func (s *DownloadSuite) SetupSuite() {
	go func() {
		if v := atomic.AddInt64(&f, 1); v == 1 {
			ctx, cancel := context.WithCancel(context.TODO())
			s.cancel = cancel
			if err := server.Start(ctx); err != nil {
				s.T().Error(err)
			}
		}
	}()
	time.Sleep(3 * time.Second) // wait start
}

func (s *DownloadSuite) TearDownSuite() {
	if s.cancel != nil {
		s.cancel()
		time.Sleep(3 * time.Second) // wait quit
	}
}

func (s *DownloadSuite) TestFileUrlList() {
	if os.Getenv("MODE") != "LOCAL" {
		s.T().Skip("not local env, skip")
	}
	req, err := http.NewRequest("GET", "http://127.0.0.1:1080/download/list", nil)
	require.Nil(s.T(), err)
	resp, err := http.DefaultClient.Do(req)
	require.Nil(s.T(), err)
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	require.Nil(s.T(), err)
	fmt.Println(string(data))
}

func (s *DownloadSuite) TestCreateUploadFile() {
	if os.Getenv("MODE") != "LOCAL" {
		s.T().Skip("not local env, skip")
	}
	client.ClientDownloadFlag.DownloadUrl = "http://localhost:1080/download"
	client.ClientDownloadFlag.FileName = "testdownload.txt"
	client.ClientDownloadFlag.DownloadPath = "./testdata/downloadstore"
	if err := client.DownloadStart(); err != nil {
		s.T().Error(err)
	}
}

func (s *DownloadSuite) TestResumeUploadFile() {
	if os.Getenv("MODE") != "LOCAL" {
		s.T().Skip("not local env, skip")
	}
	_, err := os.Stat("./testdata/video.mp4")
	if errors.Is(err, os.ErrNotExist) {
		s.T().Skip("大文件未加入版本控制中，要测试请手动加入")
	}
	client.ClientDownloadFlag.DownloadUrl = "http://localhost:1080/download"
	client.ClientDownloadFlag.FileName = "video.mp4"
	client.ClientDownloadFlag.DownloadPath = "./testdata/downloadstore"

	if err := client.DownloadStart(); err != nil {
		s.T().Error(err)
	}
}
