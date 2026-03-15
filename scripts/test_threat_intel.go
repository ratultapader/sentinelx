package main

import (
	"sentinelx/collector"
	"sentinelx/detection"
)

func main() {

	event := collector.SecurityEvent{
		EventType: "tcp_connection",
		SourceIP:  "185.220.101.12",
		Protocol:  "TCP",
	}

	detection.ThreatIntel.ProcessEvent(event)

}
