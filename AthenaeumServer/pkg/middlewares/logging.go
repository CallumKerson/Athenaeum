package middlewares

import (
	"net/http"
	"time"

	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/pkg/logging"
)

func WithRequestLogging(logger logging.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		wrappedWr := wrap(wr, logger)

		start := time.Now()
		defer logRequest(logger, start, wrappedWr, req)

		next.ServeHTTP(wrappedWr, req)
	})
}

func logRequest(logger logging.Logger, startedAt time.Time, wr *wrappedWriter, req *http.Request) {
	duration := time.Since(startedAt)

	info := map[string]interface{}{
		"latency": duration,
		"status":  wr.wroteStatus,
	}

	logger.
		WithFields(requestInfo(req)).
		WithFields(info).
		Infof("request completed with code %d", wr.wroteStatus)
}

func requestInfo(req *http.Request) map[string]interface{} {
	return map[string]interface{}{
		"path":   req.URL.Path,
		"query":  req.URL.RawQuery,
		"method": req.Method,
		"client": req.RemoteAddr,
	}
}

func wrap(wr http.ResponseWriter, logger logging.Logger) *wrappedWriter {
	return &wrappedWriter{
		ResponseWriter: wr,
		Logger:         logger,
		wroteStatus:    http.StatusOK,
	}
}

type wrappedWriter struct {
	http.ResponseWriter
	logging.Logger

	wroteStatus int
}

func (wr *wrappedWriter) WriteHeader(statusCode int) {
	wr.wroteStatus = statusCode
	wr.ResponseWriter.WriteHeader(statusCode)
}
