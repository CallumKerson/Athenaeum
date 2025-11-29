package memcache

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	cacheableResponse "github.com/CallumKerson/Athenaeum/pkg/caching/response"
)

func TestGet(t *testing.T) {
	testStore := &Store{
		sync.RWMutex{},
		2,
		time.Minute,
		map[uint64][]byte{
			14974843192121052621: (&cacheableResponse.Response{
				Value:      []byte("value 1"),
				Expiration: time.Now(),
				LastAccess: time.Now(),
			}).Bytes(),
		},
	}

	tests := []struct {
		name string
		key  uint64
		want []byte
		ok   bool
	}{
		{
			"returns right response",
			14974843192121052621,
			[]byte("value 1"),
			true,
		},
		{
			"not found",
			123,
			nil,
			false,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			b, ok := testStore.Get(testCase.key)
			assert.Equal(t, testCase.ok, ok)
			got := cacheableResponse.BytesToResponse(b).Value
			assert.Equal(t, testCase.want, got)
		})
	}
}

func TestSet(t *testing.T) {
	testStore := &Store{
		sync.RWMutex{},
		2,
		time.Minute,
		make(map[uint64][]byte),
	}

	tests := []struct {
		name     string
		key      uint64
		response *cacheableResponse.Response
	}{
		{
			"sets a response cache",
			1,
			&cacheableResponse.Response{
				Value:      []byte("value 1"),
				Expiration: time.Now().Add(1 * time.Minute),
			},
		},
		{
			"sets a response cache",
			2,
			&cacheableResponse.Response{
				Value:      []byte("value 2"),
				Expiration: time.Now().Add(1 * time.Minute),
			},
		},
		{
			"sets a response cache",
			3,
			&cacheableResponse.Response{
				Value:      []byte("value 3"),
				Expiration: time.Now().Add(1 * time.Minute),
			},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			testStore.Set(testCase.key, testCase.response.Bytes(), testCase.response.Expiration)
			assert.NotNil(t, cacheableResponse.BytesToResponse(testStore.store[testCase.key]).Value)
		})
	}

	t.Run("set is thread safe", func(t *testing.T) {
		maxSize := 2
		testStore := &Store{
			sync.RWMutex{},
			maxSize,
			time.Minute,
			make(map[uint64][]byte),
		}

		var testWaitGroup sync.WaitGroup
		for i := range 100 {
			i := uint64(i)

			testWaitGroup.Add(1)
			go func() {
				defer testWaitGroup.Done()
				testStore.Set(i, nil, time.Now().Add(1*time.Hour))
			}()
		}

		testWaitGroup.Wait()

		assert.GreaterOrEqual(t, maxSize, len(testStore.store))
	})
}

func TestRelease(t *testing.T) {
	testStore := &Store{
		sync.RWMutex{},
		2,
		time.Minute,
		map[uint64][]byte{
			14974843192121052621: (&cacheableResponse.Response{
				Expiration: time.Now().Add(1 * time.Minute),
				Value:      []byte("value 1"),
			}).Bytes(),
			14974839893586167988: (&cacheableResponse.Response{
				Expiration: time.Now(),
				Value:      []byte("value 2"),
			}).Bytes(),
			14974840993097796199: (&cacheableResponse.Response{
				Expiration: time.Now(),
				Value:      []byte("value 3"),
			}).Bytes(),
		},
	}

	tests := []struct {
		name        string
		key         uint64
		storeLength int
		wantErr     bool
	}{
		{
			"removes cached response from store",
			14974843192121052621,
			2,
			false,
		},
		{
			"removes cached response from store",
			14974839893586167988,
			1,
			false,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			testStore.Release(testCase.key)
			assert.Equal(t, testCase.storeLength, len(testStore.store))
		})
	}
}

func TestReleaseAll(t *testing.T) {
	testStore := &Store{
		sync.RWMutex{},
		2,
		time.Minute,
		map[uint64][]byte{
			14974843192121052621: (&cacheableResponse.Response{
				Expiration: time.Now().Add(1 * time.Minute),
				Value:      []byte("value 1"),
			}).Bytes(),
			14974839893586167988: (&cacheableResponse.Response{
				Expiration: time.Now(),
				Value:      []byte("value 2"),
			}).Bytes(),
			14974840993097796199: (&cacheableResponse.Response{
				Expiration: time.Now(),
				Value:      []byte("value 3"),
			}).Bytes(),
		},
	}
	assert.NotEmpty(t, testStore.store)
	testStore.ReleaseAll()
	assert.Empty(t, testStore.store)
}

func TestEvict(t *testing.T) {
	testStore := &Store{
		sync.RWMutex{},
		2,
		time.Minute,
		map[uint64][]byte{
			14974843192121052621: (&cacheableResponse.Response{
				Value:      []byte("value 1"),
				Expiration: time.Now().Add(1 * time.Minute),
				LastAccess: time.Now().Add(-1 * time.Minute),
			}).Bytes(),
			14974839893586167988: (&cacheableResponse.Response{
				Value:      []byte("value 2"),
				Expiration: time.Now().Add(1 * time.Minute),
				LastAccess: time.Now().Add(-2 * time.Minute),
			}).Bytes(),
			14974840993097796199: (&cacheableResponse.Response{
				Value:      []byte("value 3"),
				Expiration: time.Now().Add(1 * time.Minute),
				LastAccess: time.Now().Add(-3 * time.Minute),
			}).Bytes(),
		},
	}
	testStore.evict()

	_, existsInStore := testStore.store[14974840993097796199]
	assert.False(t, existsInStore)
}

func TestNewStore(t *testing.T) {
	testStore := NewStore(WithCapacity(5), WithTTL(time.Minute))
	assert.Equal(t, 5, testStore.capacity)
	assert.Equal(t, time.Minute, testStore.ttl)
}

func TestNewStore_DefaultCapacity(t *testing.T) {
	store := NewStore()
	assert.Equal(t, 100, store.capacity)
	assert.Empty(t, store.ttl)
}
