package main

import (
	"sentinelx/detection"
	"sentinelx/models"
	"sentinelx/pipeline"
	"sentinelx/storage"
	"sentinelx/metrics" 
)

func main() {

	pipeline.InitEventQueue(10000)

	// Start 5 workers
	pipeline.StartWorkerPool(5, processEvent)

	select {}

}

func processEvent(event models.SecurityEvent) {

	// ✅ 1. Save event to DB
	storage.SaveEvent(event)

	// ✅ 2. Record metrics
	// (IMPORTANT: correct order)
	metrics.RecordEvent(event.SourceIP, event.Type)

	// ✅ 3. Run detection engines
	detection.ScanDetector.ProcessEvent(event)
	detection.WAF.ProcessEvent(event)
	detection.ThreatIntel.ProcessEvent(event)

}
