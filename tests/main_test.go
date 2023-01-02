package tests

// import (
// 	"context"
// 	"fmt"
// 	"os"
// 	"testing"
// 	"time"

// 	"github.com/wwqdrh/fssync/server"
// )

// func TestMain(m *testing.M) {
// 	ctx, cancel := context.WithCancel(context.TODO())
// 	go func() {
// 		if err := server.NewFileManager().Start(ctx); err != nil {
// 			fmt.Println(err)
// 		}
// 	}()
// 	code := m.Run()
// 	cancel()
// 	time.Sleep(1 * time.Second)
// 	os.Exit(code)
// }
