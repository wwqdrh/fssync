package server

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
