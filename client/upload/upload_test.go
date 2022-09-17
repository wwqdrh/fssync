package upload

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodedMetadata(t *testing.T) {
	u, err := NewUploadFromBytes([]byte(""))
	require.Nil(t, err)
	u.Metadata["filename"] = "foobar.txt"
	enc := u.EncodedMetadata()
	assert.Equal(t, "filename Zm9vYmFyLnR4dA==", enc)
}

func TestNewUploadFromFile(t *testing.T) {
	file := fmt.Sprintf("%s/%d", os.TempDir(), time.Now().Unix())

	f, err := os.Create(file)
	assert.Nil(t, err)

	err = f.Truncate(1048576) // 1 MB
	assert.Nil(t, err)

	u, err := NewUploadFromFile(f)
	assert.Nil(t, err)
	assert.NotNil(t, u)
	assert.EqualValues(t, 1048576, u.Size())
}
