package main

import (
	"fmt"

	"sentinelx/api"
	"sentinelx/collector"
	"sentinelx/correlation"
	"sentinelx/detection"
	"sentinelx/metrics"
	"sentinelx/models"
	"sentinelx/pipeline"
	"sentinelx/ruleengine"
	"sentinelx/storage"
	"sentinelx/threatfeed"
	"sentinelx/threatintel"
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

	// Initialize database
	err = storage.InitDB()
	if err != nil {
		panic(err)
	}

	// Load detection rules
	fmt.Println("Loading detection rules...")
	err = ruleengine.LoadRules()
	if err != nil {
		panic(err)
	}

	// Initialize alert engine
	fmt.Println("Initializing alert engine...")
	detection.InitAlertEngine(AlertQueueSize)

	// Start alert processor
	go detection.StartAlertProcessor()

	// Start threat feed updater
	fmt.Println("Starting threat feed updater...")
	threatfeed.StartThreatFeedUpdater()

	// temporary test only
	threatfeed.AddTestIP("::1")
	fmt.Println("Threat feed indicators loaded:", threatfeed.Count())

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

	// Start API server
	go api.StartAPIServer()

	fmt.Println("SentinelX running")

	select {}
}

func processEvent(event models.SecurityEvent) {
	// Threat feed check
	if threatfeed.IsMalicious(event.SourceIP) {
		detection.GenerateAlert(
			"THREAT_INTEL_MATCH",
			event.SourceIP,
			"Source IP matched external threat intelligence feed",
		)
	}

	// Save event
	storage.SaveEvent(event)

	// Metrics
	metrics.RecordEvent(event.SourceIP, event.EventType)

	// Correlation tracking
	correlation.RecordEvent(event.SourceIP, event.EventType)

	// Correlation detection
	if correlation.DetectMultiStage(event.SourceIP) {
		detection.GenerateAlert(
			"MULTI_STAGE_ATTACK",
			event.SourceIP,
			"Possible coordinated attack detected",
		)
	}

	// Existing threat intel logic
	if threatintel.IsMaliciousIP(event.SourceIP) {
		detection.GenerateAlert(
			"KNOWN_MALICIOUS_IP",
			event.SourceIP,
			"Connection from known malicious IP",
		)
		return
	}

	// Rule engine processing
	rule := ruleengine.ProcessEvent(event)
	if rule != nil {
		detection.GenerateAlert(
			rule.Name,
			event.SourceIP,
			"Rule engine matched: "+rule.Name,
		)
	}

	detection.ScanDetector.ProcessEvent(event)
	detection.WAF.ProcessEvent(event)
	detection.ThreatIntel.ProcessEvent(event)
}