package tests

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/wwqdrh/fssync/client"
	"github.com/wwqdrh/fssync/server"
)

func TestUploadFile(t *testing.T) {
	server.ServerFlag.Port = ":1080"
	server.ServerFlag.Store = "./testdata/store"
	server.ServerFlag.Urlpath = "/files/"

	client.ClientFlag.Host = "http://127.0.0.1:1080/files/"
	client.ClientFlag.SpecPath = "./testdata/uploadspec"
	client.ClientFlag.Uploadfile = "./testdata/testupload.txt"

	ctx, cancel := context.WithCancel(context.TODO())

	wait := sync.WaitGroup{}
	wait.Add(2)
	go func() {
		defer wait.Done()
		if err := server.Start(ctx); err != nil {
			t.Error(err)
		}
	}()
	time.Sleep(5 * time.Second) // wait server start
	go func() {
		defer wait.Done()
		defer cancel()
		if err := client.Start(); err != nil {
			t.Error(err)
		}
	}()
	wait.Wait()
}
