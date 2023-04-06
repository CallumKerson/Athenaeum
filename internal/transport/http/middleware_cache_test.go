package http

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/CallumKerson/loggerrific/tlogger"

	cacheableResponse "github.com/CallumKerson/Athenaeum/pkg/caching/response"
)

type dummyCacheStore struct {
	sync.Mutex
	store map[uint64][]byte
}

func (s *dummyCacheStore) Get(key uint64) ([]byte, bool) {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.store[key]; ok {
		return s.store[key], true
	}
	return nil, false
}

func (s *dummyCacheStore) Set(key uint64, response []byte, expiration time.Time) {
	s.Lock()
	defer s.Unlock()
	s.store[key] = response
}

func (s *dummyCacheStore) Release(key uint64) {
	s.Lock()
	defer s.Unlock()
	delete(s.store, key)
}

func (s *dummyCacheStore) ReleaseAll() {
	s.Lock()
	defer s.Unlock()
	s.store = make(map[uint64][]byte)
}

func (s *dummyCacheStore) GetTTL() time.Duration {
	return time.Minute
}

func TestMiddleware(t *testing.T) {
	httpTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "new value %s", r.URL.Path)
	})

	testCacheStore := &dummyCacheStore{
		store: map[uint64][]byte{
			14974843192121052621: (&cacheableResponse.Response{
				Value:      []byte("old value /test-1"),
				Expiration: time.Now().Add(1 * time.Minute),
			}).Bytes(),
			14974839893586167988: (&cacheableResponse.Response{
				Value:      []byte("old value /test-2"),
				Expiration: time.Now().Add(-1 * time.Minute),
			}).Bytes(),
		},
	}

	middlewares := NewMiddlewares(tlogger.NewTLogger(t), testCacheStore)
	handler := middlewares.CachingMiddleware(httpTestHandler)

	tests := []struct {
		name     string
		url      string
		method   string
		wantBody string
	}{
		{
			"returns cached response",
			"http://foo.bar/test-1",
			http.MethodGet,
			"old value /test-1",
		},
		{
			"does not cache POSTs",
			"http://foo.bar/test-1",
			http.MethodPost,
			"new value /test-1",
		},
		{
			"cache expired",
			"http://foo.bar/test-2",
			http.MethodGet,
			"new value /test-2",
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			var request *http.Request
			var err error

			request, err = http.NewRequestWithContext(context.Background(), testCase.method, testCase.url, http.NoBody)
			if err != nil {
				t.Error(err)
				return
			}

			writer := httptest.NewRecorder()
			handler.ServeHTTP(writer, request)

			assert.Equal(t, http.StatusOK, writer.Code)
			assert.Equal(t, testCase.wantBody, writer.Body.String())
		})
	}
}

func TestGenerateKey(t *testing.T) {
	tests := []struct {
		name string
		URL  string
		want uint64
	}{
		{
			"get url checksum",
			"http://foo.bar/test-1",
			14974843192121052621,
		},
		{
			"get url 2 checksum",
			"http://foo.bar/test-2",
			14974839893586167988,
		},
		{
			"get url 3 checksum",
			"http://foo.bar/test-3",
			14974840993097796199,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.want, generateKey(testCase.URL))
		})
	}
}
