package middleware

import (
	"net/http"
	"pixelbattle/internal/pixcelbattle/metrics"
	"strconv"
	"time"
)

func Metrics(metrics metrics.Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := &statusWriter{ResponseWriter: w, status: 200}
			next.ServeHTTP(ww, r)
			duration := time.Since(start).Seconds()
			path := r.URL.Path
			method := r.Method
			status := strconv.Itoa(ww.status)

			metrics.IncHTTPRequest(method, path, status)
			metrics.ObserveHTTPDuration(method, path, duration)
		})
	}
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}
