package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// RequestsTotal counts total number of requests
	RequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_gateway_requests_total",
			Help: "Total number of requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	// RequestDuration measures request duration
	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "api_gateway_request_duration_seconds",
			Help:    "Request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// ActiveRequests tracks currently active requests
	ActiveRequests = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "api_gateway_active_requests",
			Help: "Number of active requests",
		},
	)

	// CircuitBreakerState tracks circuit breaker states
	CircuitBreakerState = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "api_gateway_circuit_breaker_state",
			Help: "Circuit breaker state (0=closed, 1=open, 2=half-open)",
		},
		[]string{"service"},
	)

	// RateLimiterActiveCount tracks active rate limiters
	RateLimiterActiveCount = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "api_gateway_rate_limiter_active_limiters",
			Help: "Number of active rate limiters in memory",
		},
	)

	// RateLimiterRejectedTotal counts rejected requests
	RateLimiterRejectedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_gateway_rate_limiter_rejected_total",
			Help: "Total number of rate limited requests",
		},
		[]string{"key_type"}, // "ip" or "tenant"
	)

	// CacheHitsTotal counts cache hits
	CacheHitsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "api_gateway_cache_hits_total",
			Help: "Total number of cache hits",
		},
	)

	// CacheMissesTotal counts cache misses
	CacheMissesTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "api_gateway_cache_misses_total",
			Help: "Total number of cache misses",
		},
	)

	// GRPCConnectionsActive tracks active gRPC connections
	GRPCConnectionsActive = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "api_gateway_grpc_connections_active",
			Help: "Number of active gRPC connections per service",
		},
		[]string{"service"},
	)

	// GRPCConnectionErrors counts connection errors
	GRPCConnectionErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_gateway_grpc_connection_errors_total",
			Help: "Total number of gRPC connection errors",
		},
		[]string{"service"},
	)
)
