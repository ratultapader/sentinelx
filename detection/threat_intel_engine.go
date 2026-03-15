package detection

import (
	"fmt"
	"sentinelx/models"
)

type ThreatIntelEngine struct{}

var ThreatIntel = &ThreatIntelEngine{}

// simple malicious IP database
var maliciousIPs = map[string]string{
	"185.220.101.12": "TOR exit node",
	"45.155.205.233": "Botnet command server",
	"103.145.13.44":  "Malware distribution server",
}

func (t *ThreatIntelEngine) ProcessEvent(event models.SecurityEvent) {

	reason, exists := maliciousIPs[event.SourceIP]

	if exists {

		fmt.Println("SECURITY ALERT")
		fmt.Println("Type: malicious_ip_detected")
		fmt.Println("Severity: CRITICAL")
		fmt.Println("Source IP:", event.SourceIP)
		fmt.Println("Description: Known malicious IP -", reason)
		fmt.Println("-----------------------------------")

	}
}
