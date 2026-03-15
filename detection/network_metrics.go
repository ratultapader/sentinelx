package detection

import (
	"fmt"
	"sync"
	"time"
)

type NetworkMetrics struct {
	mu sync.Mutex

	RequestCount    int
	ConnectionCount int
	ErrorCount      int
	TotalLatency    time.Duration
	TotalBytes      int
}

var Metrics = &NetworkMetrics{}

// RecordRequest records HTTP request metrics
func RecordRequest(latency time.Duration, bytes int) {

	Metrics.mu.Lock()
	defer Metrics.mu.Unlock()

	Metrics.RequestCount++
	Metrics.TotalLatency += latency
	Metrics.TotalBytes += bytes

	// Prometheus metrics update
	RequestCounter.Inc()
	BytesCounter.Add(float64(bytes))
}

// RecordConnection records a TCP connection
func RecordConnection() {

	Metrics.mu.Lock()
	defer Metrics.mu.Unlock()

	Metrics.ConnectionCount++

	// Prometheus metric update
	ConnectionCounter.Inc()
}

// RecordError records system errors
func RecordError() {

	Metrics.mu.Lock()
	defer Metrics.mu.Unlock()

	Metrics.ErrorCount++

	// Prometheus metric update
	ErrorCounter.Inc()
}

// StartMetricsReporter prints metrics every second
func StartMetricsReporter() {

	for {

		time.Sleep(1 * time.Second)

		Metrics.mu.Lock()

		requests := Metrics.RequestCount
		connections := Metrics.ConnectionCount
		errors := Metrics.ErrorCount
		bytes := Metrics.TotalBytes

		avgLatency := time.Duration(0)

		if Metrics.RequestCount > 0 {
			avgLatency = Metrics.TotalLatency / time.Duration(Metrics.RequestCount)
		}

		// reset counters
		Metrics.RequestCount = 0
		Metrics.ConnectionCount = 0
		Metrics.ErrorCount = 0
		Metrics.TotalLatency = 0
		Metrics.TotalBytes = 0

		Metrics.mu.Unlock()

		fmt.Println("Requests/sec:", requests)
		fmt.Println("Connections/sec:", connections)
		fmt.Println("Errors/sec:", errors)
		fmt.Println("Bytes/sec:", bytes)
		fmt.Println("Avg latency:", avgLatency)
		fmt.Println("-----------")
	}
}