package http

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/time/rate"

	"github.com/CallumKerson/loggerrific"
)

type Middlewares struct {
	loggerrific.Logger
}

func NewMiddlewares(logger loggerrific.Logger) *Middlewares {
	return &Middlewares{Logger: logger}
}

func (m *Middlewares) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.Logger.WithFields(map[string]interface{}{"method": r.Method, "path": r.URL.Path}).Infoln("Handled Request")
		next.ServeHTTP(w, r)
	})
}

func TimeoutMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
		defer cancel()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

var limiter = rate.NewLimiter(1, 3)

func SevereRateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !limiter.Allow() {
			http.Error(writer, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(writer, request)
	})
}
