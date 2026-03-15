package main

import (
	"sentinelx/collector"
	"sentinelx/detection"
)

func main() {

	for i := 1; i <= 30; i++ {

		event := collector.NewSecurityEvent("connection_open")

		event.SourceIP = "192.168.1.10"
		event.DestPort = i

		detection.ScanDetector.ProcessEvent(event)
	}
}