package proxy

import (
	"github.com/prometheus/client_golang/prometheus"
	"math"
	"reflect"

	"github.com/cloudflare/cloudflared/connection"
)

// Metrics uses connection.MetricsNamespace(aka cloudflared) as namespace and connection.TunnelSubsystem
// (tunnel) as subsystem to keep them consistent with the previous qualifier.

var (
	totalRequests = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: connection.MetricsNamespace,
			Subsystem: connection.TunnelSubsystem,
			Name:      "total_requests",
			Help:      "Amount of requests proxied through all the tunnels",
		},
	)
	concurrentRequests = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: connection.MetricsNamespace,
			Subsystem: connection.TunnelSubsystem,
			Name:      "concurrent_requests_per_tunnel",
			Help:      "Concurrent requests proxied through each tunnel",
		},
	)
	concurrentWebsocketRequests = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: connection.MetricsNamespace,
			Subsystem: connection.TunnelSubsystem,
			Name:      "concurrent_websocket_requests_per_tunnel",
			Help:      "Concurrent websocket requests proxied through each tunnel",
		},
	)
	responseByCode = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: connection.MetricsNamespace,
			Subsystem: connection.TunnelSubsystem,
			Name:      "response_by_code",
			Help:      "Count of responses by HTTP status code",
		},
		[]string{"status_code"},
	)
	earlyResponseByCode = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: connection.MetricsNamespace,
			Subsystem: connection.TunnelSubsystem,
			Name:      "early_response_by_code",
			Help:      "Count of responses by HTTP status code, incremented as soon as it is known",
		},
		[]string{"status_code"},
	)
	requestErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: connection.MetricsNamespace,
			Subsystem: connection.TunnelSubsystem,
			Name:      "request_errors",
			Help:      "Count of error proxying to origin",
		},
	)
	activeTCPSessions = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: connection.MetricsNamespace,
			Subsystem: "tcp",
			Name:      "active_sessions",
			Help:      "Concurrent count of TCP sessions that are being proxied to any origin",
		},
	)
	totalTCPSessions = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: connection.MetricsNamespace,
			Subsystem: "tcp",
			Name:      "total_sessions",
			Help:      "Total count of TCP sessions that have been proxied to any origin",
		},
	)
	connectLatency = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: connection.MetricsNamespace,
			Subsystem: "proxy",
			Name:      "connect_latency",
			Help:      "Time it takes to establish and acknowledge connections in milliseconds",
			Buckets:   []float64{1, 10, 25, 50, 100, 500, 1000, 5000},
		},
	)
	connectStreamErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: connection.MetricsNamespace,
			Subsystem: "proxy",
			Name:      "connect_streams_errors",
			Help:      "Total count of failure to establish and acknowledge connections",
		},
	)
)

func init() {
	prometheus.MustRegister(
		totalRequests,
		concurrentRequests,
		concurrentWebsocketRequests,
		responseByCode,
		earlyResponseByCode,
		requestErrors,
		activeTCPSessions,
		totalTCPSessions,
		connectLatency,
		connectStreamErrors,
	)
}

func incrementRequests() {
	totalRequests.Inc()
	concurrentRequests.Inc()
}

func incrementConcurrentWebsocketRequests() {
	concurrentWebsocketRequests.Inc()
}

func decrementConcurrentRequests() {
	concurrentRequests.Dec()
}

func decrementConcurrentWebsocketRequests() {
	concurrentWebsocketRequests.Dec()
}

func incrementTCPRequests() {
	incrementRequests()
	totalTCPSessions.Inc()
	activeTCPSessions.Inc()
}

func decrementTCPConcurrentRequests() {
	decrementConcurrentRequests()
	activeTCPSessions.Dec()
}

// Reads the current value of concurrentRequests.
// The Gauge interface does not export a method to do this, so
// use reflection. This read is not technically thread-safe
// but that's fine because we just need an approximate value.
func readConcurrentRequests() uint64 {
	val := reflect.ValueOf(concurrentRequests)
	val = val.Elem()
	f := val.FieldByName("valBits")
	bits := f.Uint()
	fVal := math.Float64frombits(bits)
	return uint64(fVal)
}

// Reads the current value of concurrentWebsocketRequests.
// The Gauge interface does not export a method to do this, so
// use reflection. This read is not technically thread-safe
// but that's fine because we just need an approximate value.
func readConcurrentWebsocketRequests() uint64 {
	val := reflect.ValueOf(concurrentWebsocketRequests)
	val = val.Elem()
	f := val.FieldByName("valBits")
	bits := f.Uint()
	fVal := math.Float64frombits(bits)
	return uint64(fVal)
}
