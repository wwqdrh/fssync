package internal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFileMd5ByContent(t *testing.T) {
	md51, err := FileMd5ByContent("./testdata/testmd5.txt")
	require.Nil(t, err)
	md52, err := FileMd5ByContent("./testdata/testmd5.txt")
	require.Nil(t, err)

	require.Equal(t, md51, md52)
}

func TestFileMd5BySpec(t *testing.T) {
	md51, err := FileMd5BySpec("./testdata/testmd5.txt")
	require.Nil(t, err)
	md52, err := FileMd5BySpec("./testdata/testmd5.txt")
	require.Nil(t, err)

	require.Equal(t, md51, md52)
}
