package server

import (
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDownloadList(t *testing.T) {
	ServerFlag.ExtraPath = "./testdata/downfile"

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/download/list", nil)
	downloadList(res, req)

	body, err := ioutil.ReadAll(res.Body)
	require.Nil(t, err)
	fmt.Println(string(body))
}

func TestDownloadSpec(t *testing.T) {
	ServerFlag.ExtraPath = "./testdata/downfile"
	ServerFlag.ExtraTruncate = 100

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/download/spec?file=a.txt", nil)
	downloadSpec(res, req)

	body, err := ioutil.ReadAll(res.Body)
	require.Nil(t, err)
	fmt.Println(string(body))
}

func TestDownloadTruncate(t *testing.T) {
	ServerFlag.ExtraPath = "./testdata/downfile"
	ServerFlag.ExtraTruncate = 100

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/download/truncate?file=a.txt&trunc=0", nil)
	downloadTruncate(res, req)

	body, err := ioutil.ReadAll(res.Body)
	require.Nil(t, err)
	fmt.Println(string(body))
}
