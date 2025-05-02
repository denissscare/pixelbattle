package middleware

import (
	"net/http"
	"pixelbattle/pkg/logger"
	"time"
)

func RequestLogger(log *logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := &responseWriter{w, http.StatusOK}
			next.ServeHTTP(rw, r)
			log.Infof("%s %s → %d (%s)",
				r.Method, r.URL.Path,
				rw.status, time.Since(start),
			)
		})
	}
}

func NoLogger(next http.Handler) http.Handler {
	return next
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
