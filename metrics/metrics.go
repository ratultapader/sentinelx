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
func RecordEvent(ip string, attackType string) {
	mu.Lock()
	defer mu.Unlock()

	TotalEvents++

	if attackType != "" {
		AttackTypes[attackType]++
	}

	if ip != "" {
		AttackerIPs[ip]++
	}
}

// Call when alert is generated
func RecordAlert() {
	mu.Lock()
	defer mu.Unlock()

	TotalAlerts++
}

// Track attack type
func RecordAttackType(attackType string) {
	mu.Lock()
	defer mu.Unlock()

	AttackTypes[attackType]++
}

// Track attacker IP
func RecordAttackerIP(ip string) {
	mu.Lock()
	defer mu.Unlock()

	AttackerIPs[ip]++
}