package detection

import (
	"sync"
	"time"

	"sentinelx/models"
)

type PortScanDetector struct {
	mu sync.Mutex

	ipPorts map[string]map[int]time.Time
}

var ScanDetector = &PortScanDetector{
	ipPorts: make(map[string]map[int]time.Time),
}

// ProcessEvent analyzes connection events and returns an alert if a scan is detected.
func (d *PortScanDetector) ProcessEvent(event models.SecurityEvent) *models.Alert {
	if event.EventType != "connection_open" {
		return nil
	}

	ip := event.SourceIP
	port := event.DestPort
	now := time.Now()

	d.mu.Lock()
	defer d.mu.Unlock()

	if _, exists := d.ipPorts[ip]; !exists {
		d.ipPorts[ip] = make(map[int]time.Time)
	}

	d.ipPorts[ip][port] = now

	d.cleanupOldEntries(ip)

	portCount := len(d.ipPorts[ip])

	if portCount > 20 {
		threatScore := 0.60
		description := "Multiple ports scanned in short time window"

		// aggressive scan threshold
		if portCount > 50 {
			threatScore = 0.78
			description = "Aggressive port scan detected across many ports in short time window"
		}

		alert := models.Alert{
			ID:          generateAlertID(),
			Timestamp:   time.Now().UTC(),
			Type:        "port_scan",
			Severity:    models.SeverityHigh,
			SourceIP:    ip,
			Description: description,
			ThreatScore: threatScore,
			Status:      models.AlertStatusNew,
			Metadata: map[string]interface{}{
				"unique_ports": portCount,
			},
		}

		delete(d.ipPorts, ip)
		return &alert
	}

	return nil
}

func (d *PortScanDetector) cleanupOldEntries(ip string) {
	window := 10 * time.Second
	now := time.Now()

	for port, timestamp := range d.ipPorts[ip] {
		if now.Sub(timestamp) > window {
			delete(d.ipPorts[ip], port)
		}
	}

	if len(d.ipPorts[ip]) == 0 {
		delete(d.ipPorts, ip)
	}
}