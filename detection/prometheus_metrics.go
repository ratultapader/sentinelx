package detection

import "github.com/prometheus/client_golang/prometheus"

var (

	RequestCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "sentinelx_requests_total",
			Help: "Total HTTP requests observed",
		},
	)

	ConnectionCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "sentinelx_connections_total",
			Help: "Total TCP connections observed",
		},
	)

	ErrorCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "sentinelx_errors_total",
			Help: "Total errors detected",
		},
	)

	BytesCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "sentinelx_bytes_total",
			Help: "Total bytes processed",
		},
	)
)

func InitPrometheusMetrics() {

	prometheus.MustRegister(RequestCounter)
	prometheus.MustRegister(ConnectionCounter)
	prometheus.MustRegister(ErrorCounter)
	prometheus.MustRegister(BytesCounter)

}