package main

import (
	"fmt"

	"sentinelx/detection"
	"sentinelx/models"
	"sentinelx/pipeline"
	"sentinelx/storage"
)

func main() {

	fmt.Println("Starting SentinelX Security Platform")

	// Initialize event logger
	err := storage.InitLogger("logs/security_events.json")
	if err != nil {
		panic(err)
	}

	// Initialize alert system
	detection.InitAlertEngine(1000)

	// Start alert processor
	go detection.StartAlertProcessor()

	// Initialize event pipeline
	pipeline.InitEventQueue(10000)

	// Start worker pool
	pipeline.StartWorkerPool(5, processEvent)

	fmt.Println("SentinelX running")

	select {}
}

// processEvent runs detection engines
func processEvent(event models.SecurityEvent) {

	detection.ScanDetector.ProcessEvent(event)
	detection.WAF.ProcessEvent(event)
	detection.ThreatIntel.ProcessEvent(event)

}
