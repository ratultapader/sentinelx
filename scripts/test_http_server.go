package main

import (
	"fmt"
	"net/http"

	"sentinelx/collector"
	"sentinelx/detection"
	"sentinelx/models"
	"sentinelx/pipeline"
)

func processEvent(event models.SecurityEvent) {

	detection.ScanDetector.ProcessEvent(event)
	detection.WAF.ProcessEvent(event)
	detection.ThreatIntel.ProcessEvent(event)

}

func main() {

	// Initialize event pipeline
	pipeline.InitEventQueue(10000)

	// Initialize alert system
	detection.InitAlertEngine(1000)

	// Start alert processor
	go detection.StartAlertProcessor()

	// Start event consumer
	go pipeline.StartEventConsumer(processEvent)

	mux := http.NewServeMux()

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Login page")
	})

	handler := collector.HTTPCollector(mux)

	fmt.Println("Test server running on :8080")

	http.ListenAndServe(":8080", handler)
}
