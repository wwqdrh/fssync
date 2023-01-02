package download

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wwqdrh/fssync/pkg/store"
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

	d, err := NewDownload("url1", "a.txt", "./testdata/downloadpath", "./testdata/downloadtemp")
	require.Nil(t, err)
	client.CreateDownload(d)

	// 判读是否存在
	_, ok = st.GetOffset(d.Fingerprint)
	require.Equal(t, true, ok)
}

func TestCreateOrResumeDownloader(t *testing.T) {
	s, err := store.NewMemoryStore()
	require.Nil(t, err)
	st, ok := s.(store.DownloadStore)
	require.Equal(t, true, ok)

	d, err := NewDownload("url1", "a.txt", "./testdata/downloadpath", "./testdata/downloadtemp")
	require.Nil(t, err)
	st.SetOffset(d.Fingerprint, 10)

	client := &DownloadClient{
		Config: &DownloadConfig{
			Store:  st,
			Resume: true,
		},
	}
	_, err = client.ResumeDownload(d)
	require.Nil(t, err)
	// 判读是否存在
	_, ok = st.GetOffset(d.Fingerprint)
	require.Equal(t, true, ok)
}
