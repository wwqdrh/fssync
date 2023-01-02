package server

import (
	"fmt"
	"io"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDownloadList(t *testing.T) {
	ServerFlag.ExtraPath = "./testdata/downfile"

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/download/list", nil)
	NewFileManager().downloadList(res, req)

	body, err := io.ReadAll(res.Body)
	require.Nil(t, err)
	fmt.Println(string(body))
}

func TestDownloadSpec(t *testing.T) {
	ServerFlag.ExtraPath = "./testdata/downfile"
	ServerFlag.ExtraTruncate = 100

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/download/spec?file=a.txt", nil)
	NewFileManager().downloadSpec(res, req)

	body, err := io.ReadAll(res.Body)
	require.Nil(t, err)
	fmt.Println(string(body))
}

func TestDownloadMd5(t *testing.T) {
	ServerFlag.ExtraPath = "./testdata/downfile"
	ServerFlag.ExtraTruncate = 100

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/download/?file=a.txt", nil)
	NewFileManager().downloadMd5(res, req)
	body, err := io.ReadAll(res.Body)
	require.Nil(t, err)
	md51 := string(body)

	res = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/download/?file=a.txt", nil)
	NewFileManager().downloadMd5(res, req)
	body, err = io.ReadAll(res.Body)
	require.Nil(t, err)
	md52 := string(body)

	require.Equal(t, md51, md52)
}

func TestDownloadTruncate(t *testing.T) {
	ServerFlag.ExtraPath = "./testdata/downfile"
	ServerFlag.ExtraTruncate = 100

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/download/truncate?file=a.txt&trunc=0", nil)
	NewFileManager().downloadTruncate(res, req)

	body, err := io.ReadAll(res.Body)
	require.Nil(t, err)
	fmt.Println(string(body))
}

func TestDownloadDelete(t *testing.T) {
	ServerFlag.ExtraPath = "./testdata/downfile"
	ServerFlag.ExtraTruncate = 100

	f, err := os.OpenFile("./testdata/downfile/temp.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	require.Nil(t, err)
	require.Nil(t, f.Close())
	defer func() {
		if _, err := os.Open("./testdata/downfile/temp.txt"); os.IsNotExist(err) {
			return
		}
		require.Nil(t, os.Remove("./testdata/downfile/temp.txt"))
	}()

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/download/delete?file=temp.txt", nil)
	NewFileManager().downloadDelete(res, req)

	body, err := io.ReadAll(res.Body)
	require.Nil(t, err)
	fmt.Println(string(body))
}
