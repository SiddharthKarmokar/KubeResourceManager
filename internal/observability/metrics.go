package observability

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// RequestCount tracks total HTTP requests.
	RequestCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "k8s_opt_request_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// RequestLatency tracks HTTP request durations.
	RequestLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "k8s_opt_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// RecommendationCount tracks the number of recommendations generated.
	RecommendationCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "k8s_opt_recommendations_total",
			Help: "Total number of recommendations generated",
		},
		[]string{"severity"},
	)

	// ErrorCount tracks the number of internal errors.
	ErrorCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "k8s_opt_errors_total",
			Help: "Total number of errors encountered",
		},
		[]string{"type"},
	)
)
