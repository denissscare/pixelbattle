package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics interface {
	IncPixelsPlaced()
	IncPixelErrors()
	IncHTTPRequest(method, path, status string)
	ObserveHTTPDuration(method, path string, duration float64)
	IncActiveConnections()
	DecActiveConnections()
}

type PrometheusMetrics struct {
	pixelsPlaced      prometheus.Counter
	pixelErrors       prometheus.Counter
	httpRequests      *prometheus.CounterVec
	httpDuration      *prometheus.HistogramVec
	activeConnections prometheus.Gauge
}

func NewPrometheusMetrics() *PrometheusMetrics {
	m := &PrometheusMetrics{
		pixelsPlaced: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "pixels_placed_total", Help: "Total pixels placed",
		}),
		pixelErrors: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "pixel_errors_total", Help: "Errors on pixel placement",
		}),
		httpRequests: prometheus.NewCounterVec(
			prometheus.CounterOpts{Name: "http_requests_total", Help: "HTTP requests count"},
			[]string{"method", "path", "status"},
		),
		httpDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{Name: "http_request_duration_seconds", Help: "HTTP request duration (s)"},
			[]string{"method", "path"},
		),
		activeConnections: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "active_connections", Help: "Active WebSocket connections",
		}),
	}
	prometheus.MustRegister(
		m.pixelsPlaced, m.pixelErrors,
		m.httpRequests, m.httpDuration, m.activeConnections,
	)
	return m
}

func (m *PrometheusMetrics) IncPixelsPlaced() { m.pixelsPlaced.Inc() }
func (m *PrometheusMetrics) IncPixelErrors()  { m.pixelErrors.Inc() }
func (m *PrometheusMetrics) IncHTTPRequest(method, path, status string) {
	m.httpRequests.WithLabelValues(method, path, status).Inc()
}
func (m *PrometheusMetrics) ObserveHTTPDuration(method, path string, duration float64) {
	m.httpDuration.WithLabelValues(method, path).Observe(duration)
}

func (m *PrometheusMetrics) IncActiveConnections() {
	m.activeConnections.Inc()
}
func (m *PrometheusMetrics) DecActiveConnections() {
	m.activeConnections.Dec()
}
