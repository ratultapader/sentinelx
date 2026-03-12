package main

import (
	"fmt"                 // Used for printing output to console
	"sentinelx/collector" // Importing our collector package where SecurityEvent is defined
)

func main() {

	// Create a new security event of type "http_request"
	// This automatically generates:
	// - unique EventID
	// - current Timestamp
	// - empty Metadata map
	event := collector.NewSecurityEvent("http_request")

	// Set the IP address where the request originated
	event.SourceIP = "192.168.1.5"

	// Set the destination IP where the request is going
	event.DestIP = "10.0.0.10"

	// Define which protocol was used
	event.Protocol = "HTTP"

	// Size of the request payload (in bytes)
	event.PayloadSize = 512

	// Convert the SecurityEvent struct into JSON format
	jsonData, err := event.ToJSON()

	// If something went wrong while converting to JSON
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	// Print the JSON event to the console
	fmt.Println(string(jsonData))
}