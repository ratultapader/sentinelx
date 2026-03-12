package main

import (
	"time"
	"sentinelx/detection"
)

func main() {

	go detection.StartMetricsReporter()

	for {

		detection.RecordRequest(50*time.Millisecond, 500)
		detection.RecordConnection()

		time.Sleep(100 * time.Millisecond)

	}
}