package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog"
	"net/http"
)

// init() gets called implicitly via an import in the main package. Metrics are registered to prometheus and
// exposed to Prometheus via a /metrics endpoint which uses the Handler() function.

func init() {
	prometheus.MustRegister(HttpRequestsTotal)
	prometheus.MustRegister(HttpDurationSeconds)
	prometheus.MustRegister(HttpErrorsTotal)
	prometheus.MustRegister(DBQueriesTotal)
	prometheus.MustRegister(CacheOperationsTotal)
	prometheus.MustRegister(StorageOperationsTotal)
	prometheus.MustRegister(ImageProcessingOperationsTotal)

	slog.Info("Init step 2 complete: metrics initialized")
}

func Handler() http.Handler {
	return promhttp.Handler()
}

var (
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests processed, partitioned by status code and method.",
		},
		[]string{"method", "status_code"},
	)

	HttpDurationSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_duration_seconds",
			Help:    "Histogram of HTTP request durations.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status_code"},
	)

	HttpErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_errors_total",
			Help: "Total number of HTTP errors by status code.",
		},
		[]string{"status_code"},
	)

	DBQueriesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_queries_total",
			Help: "Total number of database queries processed, partitioned by operation.",
		},
		[]string{"operation"},
	)

	CacheOperationsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_operations_total",
			Help: "Total number of cache operations processed, partitioned by operation.",
		},
		[]string{"operation"},
	)

	StorageOperationsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "storage_operations_total",
			Help: "Total number of storage operations processed, partitioned by operation.",
		},
		[]string{"operation"},
	)

	ImageProcessingOperationsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "image_processing_operations_total",
			Help: "Total number of image processing operations processed, partitioned by operation.",
		},
		[]string{"operation"},
	)
)
