package memcache

import "time"

type Option func(s *Store)

func WithCapacity(capacity int) Option {
	return func(a *Store) {
		a.capacity = capacity
	}
}

func WithTTL(ttl time.Duration) Option {
	return func(a *Store) {
		a.ttl = ttl
	}
}
