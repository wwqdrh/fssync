package store

type MemoryStore struct {
	m map[string]interface{}
}

func NewMemoryStore() (interface{}, error) {
	return &MemoryStore{
		make(map[string]interface{}),
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
