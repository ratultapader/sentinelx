package main

import (
	"fmt"

	"sentinelx/collector"
	"sentinelx/pipeline"
)

func main() {

	fmt.Println("Container monitor started")

	// Initialize pipeline
	pipeline.InitEventQueue(10000)

	collector.StartContainerMonitor()
}
