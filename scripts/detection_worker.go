package main

import (
	"sentinelx/detection"
	"sentinelx/metrics"
	"sentinelx/models"
	"sentinelx/pipeline"
	"sentinelx/ruleengine"
	"sentinelx/storage"
	"sentinelx/threatfeed"
)

func main() {

	// Initialize alert engine
	detection.InitAlertEngine(1000)

	// Start alert processor
	go detection.StartAlertProcessor()

	// Initialize event queue
	pipeline.InitEventQueue(10000)

	// Start 5 workers
	pipeline.StartWorkerPool(5, processEvent)

	select {}
}

func processEvent(event models.SecurityEvent) {

	// Threat feed malicious IP check
	if threatfeed.IsMalicious(event.SourceIP) {
		detection.GenerateAlert(
			"THREAT_INTEL_MATCH",
			event.SourceIP,
			"Source IP matched external threat intelligence feed",
		)
	}

	// Rule engine check
	rule := ruleengine.ProcessEvent(event)
	if rule != nil {
		detection.GenerateAlert(
			rule.Name,
			event.SourceIP,
			"Rule engine matched: "+rule.Name,
		)
	}

	// Save event to DB
	storage.SaveEvent(event)

	// Record metrics after detection
	detected := event.EventType
	if rule != nil {
		if v, ok := rule.Metadata["detected_type"].(string); ok {
			detected = v
		}
	}

	metrics.RecordEvent(event.SourceIP, event.EventType, detected)

	// Run detection engines
	detection.ScanDetector.ProcessEvent(event)
	detection.WAF.ProcessEvent(event)
	detection.ThreatIntel.ProcessEvent(event)
}
