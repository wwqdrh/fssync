package store

import (
	"fmt"
	"strconv"

	"github.com/syndtr/goleveldb/leveldb"
)

type LeveldbStore struct {
	db *leveldb.DB
}

func NewLeveldbStore(path string) (interface{}, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}

	store := &LeveldbStore{db: db}
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
	return s.db.Delete([]byte(fingerprint), nil)
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
