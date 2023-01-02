package protocol

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProtocolUrl(t *testing.T) {
	assert.Equal(t, "/download/list", PDownloadList.ServerUrl())
	assert.Equal(t, "/download/spec", PDownloadSpec.ServerUrl())
	assert.Equal(t, "/download/md5", PDownloadMd5.ServerUrl())
	assert.Equal(t, "/download/truncate", PDownloadTrucate.ServerUrl())
	assert.Equal(t, "/download/delete", PDownloadDelete.ServerUrl())
	assert.Equal(t, "/404", Unknown.ServerUrl())

	assert.Equal(t, "localhost/download/list?name=name", PDownloadList.ClientUrl("localhost", url.Values{"name": []string{"name"}}))
	assert.Equal(t, "localhost/download/spec?name=name", PDownloadSpec.ClientUrl("localhost", url.Values{"name": []string{"name"}}))
	assert.Equal(t, "localhost/download/md5?name=name", PDownloadMd5.ClientUrl("localhost", url.Values{"name": []string{"name"}}))
	assert.Equal(t, "localhost/download/truncate?name=name", PDownloadTrucate.ClientUrl("localhost", url.Values{"name": []string{"name"}}))
	assert.Equal(t, "localhost/download/delete?name=name", PDownloadDelete.ClientUrl("localhost", url.Values{"name": []string{"name"}}))
	assert.Equal(t, "localhost/404", Unknown.ClientUrl("localhost", url.Values{"name": []string{"name"}}))
}
