package memcache

import (
	"sync"
	"time"

	cacheableResponse "github.com/CallumKerson/Athenaeum/pkg/caching/response"
)

type Store struct {
	mutex    sync.RWMutex
	capacity int
	ttl      time.Duration
	store    map[uint64][]byte
}

func (s *Store) Get(key uint64) ([]byte, bool) {
	s.mutex.RLock()
	response, ok := s.store[key]
	s.mutex.RUnlock()

	if ok {
		return response, true
	}

	return nil, false
}

func (s *Store) Set(key uint64, response []byte, expiration time.Time) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, ok := s.store[key]; ok {
		s.store[key] = response
		return
	}

	if len(s.store) == s.capacity {
		s.evict()
	}

	s.store[key] = response
}

func (s *Store) Release(key uint64) {
	s.mutex.RLock()
	_, ok := s.store[key]
	s.mutex.RUnlock()

	if ok {
		s.mutex.Lock()
		delete(s.store, key)
		s.mutex.Unlock()
	}
}

func (s *Store) ReleaseAll() {
	s.mutex.Lock()
	s.store = make(map[uint64][]byte, s.capacity)
	s.mutex.Unlock()
}

func (s *Store) GetTTL() time.Duration {
	return s.ttl
}

func (s *Store) evict() {
	selectedKey := uint64(0)
	lastAccess := time.Now()

	for k, v := range s.store {
		r := cacheableResponse.BytesToResponse(v)
		if r.LastAccess.Before(lastAccess) {
			selectedKey = k
			lastAccess = r.LastAccess
		}
	}

	delete(s.store, selectedKey)
}

func NewStore(opts ...Option) *Store {
	store := &Store{}
	for _, opt := range opts {
		opt(store)
	}
	if store.capacity <= 1 {
		store.capacity = 100
	}

	store.mutex = sync.RWMutex{}
	store.store = make(map[uint64][]byte, store.capacity)
	return store
}
