package store

import (
	"errors"
	"sync"

	"github.com/wwqdrh/gokit/logger"
)

var (
	_ DownloadStore = new(MemoryStore)
	_ UploadStore   = new(MemoryStore)
)

type MemoryStore struct {
	m map[string]interface{}

	combiled     map[string]bool
	maxOffset    map[string]int64 // 记录最大长度
	chunckStatus map[string][]int // 1已分配但是未下载 0未下载 2已下载 获取chunck的状态 false未下载，true已下载
	mutex        *sync.RWMutex
}

func NewMemoryStore() (interface{}, error) {
	return &MemoryStore{
		m:            make(map[string]interface{}),
		combiled:     map[string]bool{},
		maxOffset:    map[string]int64{},
		chunckStatus: map[string][]int{},
		mutex:        &sync.RWMutex{},
	}, nil
}

func (s *MemoryStore) Get(fingerprint string) (string, bool) {
	url, ok := s.m[fingerprint]
	return url.(string), ok
}

func (s *MemoryStore) Set(fingerprint, url string) error {
	s.m[fingerprint] = url
	return nil
}

func (s *MemoryStore) Delete(fingerprint string) error {
	delete(s.m, fingerprint)
	return nil
}

func (s *MemoryStore) Close() error {
	for k := range s.m {
		delete(s.m, k)
	}
	return nil
}

////////////////////
// for download
////////////////////

func (s *MemoryStore) GetOffset(fingerprint string) (int64, bool) {
	url, ok := s.m[fingerprint]
	return url.(int64), ok
}

func (s *MemoryStore) SetOffset(fingerprint string, offset int64) error {
	s.m[fingerprint] = offset
	return nil
}

// 设置一个figerprint最大的切片数
func (s *MemoryStore) SetMaxOffset(figerprint string, offset int64) error {
	s.maxOffset[figerprint] = offset
	return nil
}

func (s *MemoryStore) GetMaxOffset(figerprint string) (int64, bool) {
	val, ok := s.maxOffset[figerprint]
	return val, ok
}

// 并发安全 获取一个还未下载的切片
func (s *MemoryStore) GetBlankOffset(figerprint string) (int64, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	data, ok := s.chunckStatus[figerprint]
	if !ok {
		if v, ok := s.GetMaxOffset(figerprint); !ok {
			logger.DefaultLogger.Warn("未设置maxoffset")
			return -1, false
		} else {
			data = make([]int, v)
			s.chunckStatus[figerprint] = data
		}
	}
	for i, n := 0, len(data); i < n; i++ {
		switch data[i] {
		case 0:
			data[i] = 1 // 已分配
			return int64(i), true
		case 1, 2:
			continue
		}
	}
	return -1, false
}

// 并发安全 标记一个切片已经下载完成
func (s *MemoryStore) SetOkOffset(figerprint string, offset int64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	data, ok := s.chunckStatus[figerprint]
	if !ok {
		if v, ok := s.GetMaxOffset(figerprint); !ok {
			// logger.DefaultLogger.Warn("未设置maxoffset")
			return errors.New("未设置maxoffset")
		} else {
			data = make([]int, v)
			s.chunckStatus[figerprint] = data
		}
	}
	if offset > int64(len(data)) {
		return errors.New("offset illegal")
	}
	data[offset] = 2
	return nil
}

// 并发安全 标记一个切片已经下载失败
func (s *MemoryStore) SetFailOffset(figerprint string, offset int64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	data, ok := s.chunckStatus[figerprint]
	if !ok {
		if v, ok := s.GetMaxOffset(figerprint); !ok {
			// logger.DefaultLogger.Warn("未设置maxoffset")
			return errors.New("未设置maxoffset")
		} else {
			data = make([]int, v)
			s.chunckStatus[figerprint] = data
		}
	}
	if offset > int64(len(data)) {
		return errors.New("offset illegal")
	}
	data[offset] = 0
	return nil
}

func (s *MemoryStore) IsDone(figerprint string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	val, ok := s.chunckStatus[figerprint]
	if !ok {
		return false
	}
	for _, item := range val {
		if item != 2 {
			return false
		}
	}
	return true
}

func (s *MemoryStore) IsCombile(figerprint string) error {
	if !s.combiled[figerprint] {
		return errors.New("未合并")
	}
	return nil
}

func (s *MemoryStore) SetCombile(figerprint string) error {
	s.combiled[figerprint] = true
	return nil
}
