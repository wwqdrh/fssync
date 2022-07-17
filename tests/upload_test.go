package tests

import (
	"context"
	"errors"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/wwqdrh/fssync/client"
	"github.com/wwqdrh/fssync/server"
)

var f = int64(0) // atomic

func init() {
	server.ServerFlag.Port = ":1080"
	server.ServerFlag.Store = "./testdata/store"
	server.ServerFlag.Urlpath = "/files/"
	server.ServerFlag.Store = "./testdata/store"

	client.ClientUploadFlag.Host = "http://127.0.0.1:1080/files/"
	client.ClientUploadFlag.SpecPath = "./testdata/uploadspec"
}

type UploadSuite struct {
	suite.Suite

	cancel context.CancelFunc
}

func TestUploadSuite(t *testing.T) {
	suite.Run(t, &UploadSuite{})
}

func (s *UploadSuite) SetupSuite() {
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

func (s *UploadSuite) TearDownSuite() {
	if s.cancel != nil {
		s.cancel()
		time.Sleep(3 * time.Second) // wait quit
	}
}

func (s *UploadSuite) TestCreateUploadFile() {
	if os.Getenv("MODE") != "LOCAL" {
		s.T().Skip("not local env, skip")
	}
	client.ClientUploadFlag.Uploadfile = "./testdata/testupload.txt"
	if err := client.UploadStart(); err != nil {
		s.T().Error(err)
	}
}

func (s *UploadSuite) TestResumeUploadFile() {
	if os.Getenv("MODE") != "LOCAL" {
		s.T().Skip("not local env, skip")
	}
	_, err := os.Stat("./testdata/video.mp4")
	if errors.Is(err, os.ErrNotExist) {
		s.T().Skip("大文件未加入版本控制中，要测试请手动加入")
	}
	client.ClientUploadFlag.Uploadfile = "./testdata/video.mp4"

	if err := client.UploadStart(); err != nil {
		s.T().Error(err)
	}
}
