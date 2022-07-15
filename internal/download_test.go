package internal

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wwqdrh/fssync/internal/store"
)

func TestCreateDownloader(t *testing.T) {
	s, err := store.NewMemoryStore()
	require.Nil(t, err)
	st, ok := s.(store.DownloadStore)
	require.Equal(t, true, ok)

	client := &DownloadClient{
		Config: &DownloadConfig{
			Store: st,
		},
	}

	d, err := NewDownload("url1")
	require.Nil(t, err)
	client.CreateDownload(d)

	// 判读是否存在
	_, ok = st.GetOffset("url1")
	require.Equal(t, true, ok)
}

func TestCreateOrResumeDownloader(t *testing.T) {
	s, err := store.NewMemoryStore()
	require.Nil(t, err)
	st, ok := s.(store.DownloadStore)
	require.Equal(t, true, ok)
	st.SetOffset("url1", 10)

	client := &DownloadClient{
		Config: &DownloadConfig{
			Store:  st,
			Resume: true,
		},
	}

	d, err := NewDownload("url1")
	require.Nil(t, err)
	_, err = client.ResumeDownload(d)
	require.Nil(t, err)
	// 判读是否存在
	_, ok = st.GetOffset("url1")
	require.Equal(t, true, ok)
}
