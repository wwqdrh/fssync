package server

import (
	"reflect"
	"sort"
	"testing"

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
