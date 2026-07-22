package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	CacheHits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		},
		[]string{"cache_key"},
	)
	CacheMisses = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		},
		[]string{"cache_key"},
	)
	RedisLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "redis_latency_seconds",
			Help:    "Latency of Redis operations",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)
	QueueLength = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "queue_length",
			Help: "Current length of the background job queue",
		},
		[]string{"queue_name"},
	)
	WorkerJobsProcessed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "worker_jobs_processed_total",
			Help: "Total number of jobs processed by worker",
		},
		[]string{"queue_name", "status"},
	)
	QueueProcessingDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "worker_job_processing_duration_seconds",
			Help:    "Time taken to process a background job",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"queue_name"},
	)
	UploadsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "spaces_uploads_total",
			Help: "Total number of file uploads to Spaces",
		},
		[]string{"status"}, // "success", "failure"
	)
	UploadLatency = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "spaces_upload_duration_seconds",
			Help:    "Latency of object uploads to Spaces",
			Buckets: prometheus.DefBuckets,
		},
	)
	TaskCreationDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "task_creation_duration_seconds",
			Help:    "Time taken to create a task (including DB and queue ops)",
			Buckets: prometheus.DefBuckets,
		},
	)
)
