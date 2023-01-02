package store

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLevelDB(t *testing.T) {
	l, err := NewLeveldbStore("./testdata/leveldb")
	require.Nil(t, err)

	lv, ok := l.(DownloadStore)
	require.Equal(t, true, ok)

	figure := "testfigure"
	_, ok = lv.GetBlankOffset(figure)
	require.Equal(t, false, ok)

	lv.SetMaxOffset(figure, 10)
	wait := sync.WaitGroup{}
	wait.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wait.Done()
			val, ok := lv.GetBlankOffset(figure)
			assert.Equal(t, true, ok)
			fmt.Println(val)
		}()
	}
	wait.Wait()
}
