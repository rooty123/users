package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// RequestCounter tracks total number of requests since application start
	RequestCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "users_api_requests_total",
			Help: "Total number of requests since application start",
		},
		[]string{"method", "endpoint", "status"},
	)

	// RequestDuration tracks request duration in seconds
	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "users_api_request_duration_seconds",
			Help:    "Request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// UsersGauge tracks total number of users in the system
	UsersGauge = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "users_api_total_users",
			Help: "Total number of users in the system",
		},
	)
)
