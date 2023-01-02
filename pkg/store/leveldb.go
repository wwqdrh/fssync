package store

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

var (
	_ DownloadStore = new(LeveldbStore)
	_ UploadStore   = new(LeveldbStore)
)

type LeveldbStore struct {
	db *leveldb.DB

	chunkStatus map[string][]int
	mutex       *sync.RWMutex
}

func NewLeveldbStore(path string) (interface{}, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}

	store := &LeveldbStore{db: db,
		chunkStatus: map[string][]int{},
		mutex:       &sync.RWMutex{},
	}
	return store, err
}

func (s *LeveldbStore) Get(fingerprint string) (string, bool) {
	url, err := s.db.Get([]byte(fingerprint), nil)
	ok := true
	if err != nil {
		ok = false
	}
	return string(url), ok
}

func (s *LeveldbStore) Set(fingerprint, url string) error {
	return s.db.Put([]byte(fingerprint), []byte(url), nil)
}

func (s *LeveldbStore) Delete(fingerprint string) error {
	items := []string{fingerprint, fingerprint + "maxoffset", fingerprint + "chunkstatus", fingerprint + "combiled"}
	for _, item := range items {
		if err := s.db.Delete([]byte(item), nil); err != nil {
			return err
		}
	}
	return nil
}

func (s *LeveldbStore) Close() error {
	return s.db.Close()
}

////////////////////
// for download
////////////////////

func (s *LeveldbStore) GetOffset(fingerprint string) (int64, bool) {
	val, err := s.db.Get([]byte(fingerprint), nil)
	if err != nil {
		return -1, false
	}

	v, err := strconv.ParseInt(string(val), 10, 64)
	if err != nil {
		return -1, false
	}
	return v, true
}

func (s *LeveldbStore) SetOffset(fingerprint string, offset int64) error {
	return s.db.Put([]byte(fingerprint), []byte(fmt.Sprint(offset)), nil)
}

// 设置一个figerprint最大的切片数
func (s *LeveldbStore) SetMaxOffset(figerprint string, offset int64) error {
	if err := s.db.Put([]byte(figerprint+"maxoffset"), []byte(fmt.Sprint(offset)), nil); err != nil {
		return err
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.updateStatus(figerprint)
}

func (s *LeveldbStore) GetMaxOffset(figerprint string) (int64, bool) {
	val, err := s.db.Get([]byte(figerprint+"maxoffset"), nil)
	if err != nil {
		return -1, false
	}

	v, err := strconv.ParseInt(string(val), 10, 64)
	if err != nil {
		return -1, false
	}
	return v, true
}

// 并发安全 获取一个还未下载的切片
// leveldb val : 1,1,1,1,1,1,1
func (s *LeveldbStore) GetBlankOffset(figerprint string) (int64, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	val, ok := s.chunkStatus[figerprint]
	if !ok {
		if err := s.updateStatus(figerprint); err != nil {
			return -1, false
		}
	}

	for i, n := 0, len(val); i < n; i++ {
		switch val[i] {
		case 0:
			val[i] = 1 // 已分配
			return int64(i), true
		case 1, 2:
			continue
		}
	}
	return -1, false
}

// 并发安全 标记一个切片已经下载完成
func (s *LeveldbStore) SetOkOffset(figerprint string, offset int64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	val, ok := s.chunkStatus[figerprint]
	if !ok {
		if err := s.updateStatus(figerprint); err != nil {
			return err
		}
	}

	val[offset] = 2
	return s.flushStatus(figerprint)
}

func (s *LeveldbStore) SetFailOffset(figerprint string, offset int64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	val, ok := s.chunkStatus[figerprint]
	if !ok {
		if err := s.updateStatus(figerprint); err != nil {
			return err
		}
	}

	val[offset] = 0
	return s.flushStatus(figerprint)
}

func (s *LeveldbStore) updateStatus(figerprint string) error {
	valByte, err := s.db.Get([]byte(figerprint+"chunkstatus"), nil)
	if err != nil {
		maxOff, ok := s.GetMaxOffset(figerprint)
		if !ok {
			return errors.New("未设置最大offset")
		}
		s.chunkStatus[figerprint] = make([]int, maxOff)
		if err := s.flushStatus(figerprint); err != nil {
			return err
		}
		return nil
	}

	statusStr := strings.Split(string(valByte), ",")
	status := make([]int, len(statusStr))
	for i, item := range statusStr {
		switch item {
		case "0":
			status[i] = 0
		case "1":
			status[i] = 1
		case "2":
			status[i] = 2
		default:
			status[i] = 0
		}
	}
	s.chunkStatus[figerprint] = status
	return nil
}

func (s *LeveldbStore) flushStatus(figerprint string) error {
	val, ok := s.chunkStatus[figerprint]
	if !ok {
		return errors.New("figerprint不存在")
	}

	statusStr := make([]string, len(val))
	for i, item := range val {
		statusStr[i] = fmt.Sprint(item)
	}

	return s.db.Put([]byte(figerprint+"chunkstatus"), []byte(strings.Join(statusStr, ",")), nil)
}

func (s *LeveldbStore) IsDone(figerprint string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	val, ok := s.chunkStatus[figerprint]
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

func (s *LeveldbStore) IsCombile(figerprint string) error {
	valByte, err := s.db.Get([]byte(figerprint+"combiled"), nil)
	if err != nil {
		return err
	}
	if string(valByte) != "1" {
		return errors.New("未合并")
	}
	return nil
}

func (s *LeveldbStore) SetCombile(figerprint string) error {
	return s.db.Put([]byte(figerprint+"combiled"), []byte(fmt.Sprint(1)), nil)
}
