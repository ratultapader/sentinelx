package detection

import (
	"fmt"
	"time"

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

// ProcessEvent checks whether the source IP matches known threat intelligence.
func (t *ThreatIntelEngine) ProcessEvent(event models.SecurityEvent) *models.Alert {

	reason, exists := maliciousIPs[event.SourceIP]
	if !exists {
		return nil
	}

	// ✅ CREATE ALERT WITH TENANT
	alert := models.Alert{
		ID:          generateAlertID(),
		TenantID:    event.TenantID, // 🔥 FIXED
		Timestamp:   time.Now().UTC(),
		Type:        "threat_intel_match",
		Severity:    models.SeverityCritical,
		SourceIP:    event.SourceIP,
		Description: "Source IP matched known malicious threat feed",
		ThreatScore: 0.98,
		Status:      models.AlertStatusNew,
		Metadata: map[string]interface{}{
			"matched_ip": event.SourceIP,
			"reason":     reason,
		},
	}

	// ✅ DEBUG (OPTIONAL)
	fmt.Println("DEBUG ALERT TENANT:", alert.TenantID)

	return &alert
}