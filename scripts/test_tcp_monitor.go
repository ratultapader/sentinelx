package main

import (
	"fmt"                 // Used for printing messages to the console
	"sentinelx/collector" // Import the collector package where TCP monitor is implemented
)

func main() {

	// Print a message so we know the TCP monitor has started
	fmt.Println("TCP monitor running on port 9000")

	// Start the TCP monitoring server on port 9000
	// This function will:
	// 1. Open a TCP listener on port 9000
	// 2. Accept incoming connections
	// 3. Track connection open/close events
	// 4. Measure bytes transferred and connection duration
	// 5. Generate SecurityEvents for monitoring
	collector.StartTCPMonitor("9000")

}