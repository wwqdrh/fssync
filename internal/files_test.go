package internal

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	_, err := os.Stat("./testdata/testmd5.txt")
	if os.IsNotExist(err) {
		if _, err := os.Stat("./testdata"); os.IsNotExist(err) {
			if err := os.Mkdir("testdata", os.ModePerm); err != nil {
				fmt.Print(err)
			}
		}
	}

	f, err := os.Create("./testdata/testmd5.txt")
	if err != nil {
		fmt.Println(err)
	}
	f.WriteString("hello testmd5.txt")
	f.Close()
	m.Run()
}

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
