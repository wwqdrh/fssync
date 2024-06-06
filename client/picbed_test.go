package client

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(t *testing.M) {
	godotenv.Load("./.env")

	os.Exit(t.Run())
}

func TestParsePid(t *testing.T) {
	weiboResp := `
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
    <script type="text/javascript">document.domain="sina.com.cn";</script>
{"code":"A00006","data":{"count":1,"data":"eyJ1aWQiOjU0ODc5MjIxNjQsImFwcCI6Im1pbmlibG9nIiwiY291bnQiOjEsInRpbWUiOjE3MTc2MDA0NzYuNjYzLCJwaWNzIjp7InBpY18xIjp7IndpZHRoIjo0ODYsInNpemUiOjE4OTkzLCJyZXQiOjEsImhlaWdodCI6NDc4LCJuYW1lIjoicGljXzEiLCJwaWQiOiIwMDVab0xmQ2d5MWhxZXZ3MXI4Z2FqMzBkaTBkYWpybSJ9fX0=","pics":{"pic_1":{"width":486,"size":18993,"ret":1,"height":478,"name":"pic_1","pid":"005ZoLfCgy1hqevw1r8gaj30di0dajrm"}}}}
	`

	id, err := parsePicPid(weiboResp)
	if err != nil {
		t.Error(err)
		return
	}
	if id != "005ZoLfCgy1hqevw1r8gaj30di0dajrm" {
		t.Error("解析失败")
	}
}

func TestWeiboUploader_Upload(t *testing.T) {
	uploader := NewWeiboUploader()
	fileURL := "testdata/pic.jpg" // 替换为您的测试文件路径
	username := ""
	password := ""
	cookie := os.Getenv("COOKIE")
	domain := "https://ww1.sinaimg.cn"
	quality := "large"
	cookieMode := true

	uploadURL, _, err := uploader.Upload(fileURL, username, password, cookie, domain, quality, cookieMode)
	if err != nil {
		t.Errorf("Upload failed: %s", err)
		return
	}

	t.Logf("Upload successful, URL: %s", uploadURL)
}
