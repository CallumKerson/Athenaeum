package http

import "github.com/CallumKerson/loggerrific"

type HandlerOption func(h *Handler)

func WithVersion(version string) HandlerOption {
	return func(h *Handler) {
		h.version = version
	}
}

func WithLogger(logger loggerrific.Logger) HandlerOption {
	return func(h *Handler) {
		h.Log = logger
	}
}

func WithCacheStore(cacheStore CacheStore) HandlerOption {
	return func(h *Handler) {
		h.CacheStore = cacheStore
	}
}
