package main

import (
	"fmt"
	"net/http"

	// Import SentinelX detection package where Prometheus metrics are defined
	"sentinelx/detection"

	// Prometheus HTTP handler used to expose metrics to Prometheus server
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	// Initialize and register all Prometheus metrics
	// This must be done before starting the metrics server
	detection.InitPrometheusMetrics()

	// Register HTTP endpoint "/metrics"
	// Prometheus will scrape this endpoint to collect metrics
	http.Handle("/metrics", promhttp.Handler())

	// Print message to terminal indicating the metrics server is running
	fmt.Println("Metrics server running on :9090")

	// Start HTTP server on port 9090
	// This server exposes the /metrics endpoint for Prometheus
	http.ListenAndServe(":9090", nil)

}