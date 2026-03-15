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

func (d *PortScanDetector) ProcessEvent(event models.SecurityEvent) {

	if event.EventType != "connection_open" {
		return
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

	if len(d.ipPorts[ip]) > 20 {

		GenerateAlert(
	"HIGH",
	"port_scan",
	ip,
	"Multiple ports scanned in short time",
)

		delete(d.ipPorts, ip)
	}
}

func (d *PortScanDetector) cleanupOldEntries(ip string) {

	window := 10 * time.Second
	now := time.Now()

	for port, timestamp := range d.ipPorts[ip] {

		if now.Sub(timestamp) > window {
			delete(d.ipPorts[ip], port)
		}
	}
}