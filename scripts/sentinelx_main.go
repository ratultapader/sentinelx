package main

import (
	"fmt"

	"sentinelx/collector"
	"sentinelx/detection"
	"sentinelx/models"
	"sentinelx/pipeline"
	"sentinelx/storage"
	"sentinelx/threatintel"
	"sentinelx/metrics"
	"sentinelx/api"
	
	
)

const (
	EventQueueSize = 10000
	AlertQueueSize = 1000
	WorkerCount    = 5
)

func main() {

	fmt.Println("Starting SentinelX Security Platform")

	// Load threat intelligence feed
	err := threatintel.LoadThreatFeed("data/malicious_ips.txt")
	if err != nil {
		panic(err)
	}

	// Initialize logger
	fmt.Println("Initializing event logger...")
	err = storage.InitLogger("logs/security_events.json")
	if err != nil {
		panic(err)
	}

	// Initialize alert engine
	fmt.Println("Initializing alert engine...")
	detection.InitAlertEngine(AlertQueueSize)

	// Start alert processor
	go detection.StartAlertProcessor()

	// Initialize event pipeline
	fmt.Println("Initializing event pipeline...")
	pipeline.InitEventQueue(EventQueueSize)

	// Start worker pool
	fmt.Println("Starting worker pool...")
	pipeline.StartWorkerPool(WorkerCount, processEvent)

	// Start metrics reporter
go metrics.StartMetricsReporter()

	// Start HTTP collector
	go collector.StartHTTPServer()

	go api.StartAPIServer()

	fmt.Println("SentinelX running")

	select {}
}

func processEvent(event models.SecurityEvent) {

	metrics.RecordEvent(event.EventType, event.SourceIP)

	if threatintel.IsMaliciousIP(event.SourceIP) {

		detection.GenerateAlert(
			"KNOWN_MALICIOUS_IP",
			event.SourceIP,
			"Connection from known malicious IP",
		)

		return
	}

	detection.ScanDetector.ProcessEvent(event)
	detection.WAF.ProcessEvent(event)
	detection.ThreatIntel.ProcessEvent(event)
}