package server

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListDirFile(t *testing.T) {
	source := "./testdata/files"

	sourceRes, err := ListDirFile(source, false)
	require.Nil(t, err)
	sort.Slice(sourceRes, func(i, j int) bool { return sourceRes[i] < sourceRes[j] })

	require.Equal(t,
		true,
		reflect.DeepEqual(
			sourceRes,
			[]string{"a.txt", "a/a.txt", "b.txt"},
		))
}

func TestGetFileSpecInfo(t *testing.T) {
	file := fmt.Sprintf("%s/%d", os.TempDir(), time.Now().Unix())

	f, err := os.Create(file)
	assert.Nil(t, err)
	defer os.Remove(f.Name())

	err = f.Truncate(10000)
	assert.Nil(t, err)

	res, err := GetFileSpecInfo(f.Name(), 1000)
	assert.Nil(t, err)
	assert.EqualValues(t, 10, res)

	err = f.Truncate(10001)
	assert.Nil(t, err)

	res, err = GetFileSpecInfo(f.Name(), 1000)
	assert.Nil(t, err)
	assert.EqualValues(t, 11, res)
}

func TestGetFileData(t *testing.T) {
	file := fmt.Sprintf("%s/%d", os.TempDir(), time.Now().Unix())

	f, err := os.Create(file)
	assert.Nil(t, err)
	defer os.Remove(f.Name())

	err = f.Truncate(10000)
	assert.Nil(t, err)

	res, err := GetFileData(f.Name(), 10, 100)
	require.Nil(t, err)
	require.EqualValues(t, 100, len(res))
}

func TestIsSubDir(t *testing.T) {
	tables := []struct {
		cur    string
		sub    string
		expect bool
	}{
		{"testdata", "testdata/a.txt", true},
		{"testdata", "testdata/../a.txt", false},
		{"testdata", "testdata/a/../../a.txt", false},
		{"testdata", "testdata/a/b/../a.txt", true},
	}

	for i, item := range tables {
		if item.expect != isSubDir(item.cur, item.sub) {
			t.Error(fmt.Sprintf("%d - 错误", i))
		}
	}
}
