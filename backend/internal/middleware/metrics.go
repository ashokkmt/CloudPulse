package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

type responseWriterObserver struct {
	http.ResponseWriter
	status int
}

func (w *responseWriterObserver) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriterObserver) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	return w.ResponseWriter.Write(b)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		path := r.URL.Path
		// Generalize task ids to prevent metric cardinality explosion
		if len(path) > 11 && path[:11] == "/api/tasks/" {
			path = "/api/tasks/:id"
		}
		
		observer := &responseWriterObserver{ResponseWriter: w, status: http.StatusOK}
		
		next.ServeHTTP(observer, r)
		
		duration := time.Since(start).Seconds()
		statusStr := strconv.Itoa(observer.status)
		
		httpRequestsTotal.WithLabelValues(r.Method, path, statusStr).Inc()
		httpRequestDuration.WithLabelValues(r.Method, path).Observe(duration)
	})
}
