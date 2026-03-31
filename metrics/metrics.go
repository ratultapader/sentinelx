package metrics

import "sync"

// Global mutex (shared across package)
var mu sync.RWMutex

// Core counters
var TotalEvents int
var TotalAlerts int

// Aggregation maps
var AttackTypes = make(map[string]int)
var AttackerIPs = make(map[string]int)

//
// =====================
// EVENT METRICS
// =====================
//

// Call when event is processed
func RecordEvent(ip string, eventType string, detectedType string) {
	mu.Lock()
	defer mu.Unlock()

	TotalEvents++

	if detectedType != "" {
		AttackTypes[detectedType]++
	} else if eventType != "" {
		AttackTypes[eventType]++
	}

	if ip != "" {
		AttackerIPs[ip]++
	}
}

// Call when alert is generated
func RecordAlert(attackType string) {
	mu.Lock()
	defer mu.Unlock()

	TotalAlerts++

	if attackType != "" {
		AttackTypes[attackType]++
	}
}
