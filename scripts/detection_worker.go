package main

import (
	"sentinelx/detection"
	"sentinelx/models"
	"sentinelx/pipeline"
)

func main() {

	pipeline.InitEventQueue(10000)

	// Start 5 workers
	pipeline.StartWorkerPool(5, processEvent)

	select {}

}

func processEvent(event models.SecurityEvent) {

	detection.ScanDetector.ProcessEvent(event)
	detection.WAF.ProcessEvent(event)
	detection.ThreatIntel.ProcessEvent(event)

}
