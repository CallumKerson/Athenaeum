package http

import (
	"hash/fnv"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	cacheableResponse "github.com/CallumKerson/Athenaeum/pkg/caching/response"
)

func (m *Middlewares) CachingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if (m.CacheStore != nil) && (request.Method == http.MethodGet) {
			key := generateKey(request.URL.String())

			b, ok := m.CacheStore.Get(key)
			response := cacheableResponse.BytesToResponse(b)
			if ok {
				if response.Expiration.After(time.Now()) {
					response.LastAccess = time.Now()
					m.CacheStore.Set(key, response.Bytes(), response.Expiration)

					for k, v := range response.Header {
						writer.Header().Set(k, strings.Join(v, ","))
					}
					_, _ = writer.Write(response.Value)
					return
				}

				m.CacheStore.Release(key)
			}

			rec := httptest.NewRecorder()
			next.ServeHTTP(rec, request)
			result := rec.Result()
			defer result.Body.Close()

			statusCode := result.StatusCode
			value := rec.Body.Bytes()
			if statusCode < 400 {
				now := time.Now()

				response := &cacheableResponse.Response{
					Value:      value,
					Header:     result.Header,
					Expiration: now.Add(m.CacheStore.GetTTL()),
					LastAccess: now,
				}
				m.CacheStore.Set(key, response.Bytes(), response.Expiration)
			}
			for k, v := range result.Header {
				writer.Header().Set(k, strings.Join(v, ","))
			}
			writer.WriteHeader(statusCode)
			_, _ = writer.Write(value)
			return
		}
		next.ServeHTTP(writer, request)
	})
}

func generateKey(keyURL string) uint64 {
	hash := fnv.New64a()
	hash.Write([]byte(keyURL))

	return hash.Sum64()
}
