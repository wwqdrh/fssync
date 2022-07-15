package tests

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/wwqdrh/fssync/client"
	"github.com/wwqdrh/fssync/server"
)

type UploadSuite struct {
	suite.Suite

	cancel context.CancelFunc
}

func TestUploadSuite(t *testing.T) {

	suite.Run(t, &UploadSuite{})
}

func (s *UploadSuite) SetupSuite() {
	server.ServerFlag.Port = ":1080"
	server.ServerFlag.Store = "./testdata/store"
	server.ServerFlag.Urlpath = "/files/"

	client.ClientFlag.Host = "http://127.0.0.1:1080/files/"
	client.ClientFlag.SpecPath = "./testdata/uploadspec"

	ctx, cancel := context.WithCancel(context.TODO())
	s.cancel = cancel
	go func() {
		if err := server.Start(ctx); err != nil {
			s.T().Error(err)
		}
	}()
}

func (s *UploadSuite) TearDownSuite() {
	s.cancel()
	time.Sleep(3 * time.Second) // wait quit
}

func (s *UploadSuite) TestCreateUploadFile() {
	client.ClientFlag.Uploadfile = "./testdata/testupload.txt"
	if err := client.Start(); err != nil {
		s.T().Error(err)
	}
}

func (s *UploadSuite) TestResumeUploadFile() {
	_, err := os.Stat("./testdata/video.mp4")
	if errors.Is(err, os.ErrNotExist) {
		s.T().Skip("大文件未加入版本控制中，要测试请手动加入")
	}
	client.ClientFlag.Uploadfile = "./testdata/video.mp4"

	if err := client.Start(); err != nil {
		s.T().Error(err)
	}
}
