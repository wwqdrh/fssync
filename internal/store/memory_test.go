package store

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryStore(t *testing.T) {
	s, err := NewMemoryStore()
	require.Nil(t, err)

	sv, ok := s.(DownloadStore)
	require.Equal(t, true, ok)

	figure := "testoffset"
	_, ok = sv.GetBlankOffset(figure)
	require.Equal(t, false, ok)

	sv.SetMaxOffset(figure, 10)
	wait := sync.WaitGroup{}
	wait.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wait.Done()
			val, ok := sv.GetBlankOffset(figure)
			assert.Equal(t, true, ok)
			fmt.Println(val)
		}()
	}
	wait.Wait()
}
