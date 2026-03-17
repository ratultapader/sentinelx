package correlation

import (
	"sync"
	"time"
)

// =====================
// DATA STRUCTURE
// =====================

// attackRecord stores recent attack activity per IP
type attackRecord struct {
	events    []string
	lastSeen  time.Time
	triggered bool // prevents duplicate alerts
}

var attackCache = make(map[string]*attackRecord)
var mu sync.Mutex

// =====================
// RECORD EVENTS
// =====================

// RecordEvent stores attack activity for correlation
func RecordEvent(ip string, eventType string) {

	mu.Lock()
	defer mu.Unlock()

	record, exists := attackCache[ip]

	// If new IP → create record
	if !exists {
		attackCache[ip] = &attackRecord{
			events:    []string{eventType},
			lastSeen:  time.Now(),
			triggered: false,
		}
		return
	}

	// 🔥 Reset window after 30 seconds (NEW ATTACK WINDOW)
	if time.Since(record.lastSeen) > 30*time.Second {
		record.events = []string{eventType}
		record.triggered = false
		record.lastSeen = time.Now()
		return
	}

	// Append event
	record.events = append(record.events, eventType)
	record.lastSeen = time.Now()

	// Keep only last 5 events (memory control)
	if len(record.events) > 5 {
		record.events = record.events[1:]
	}
}

// =====================
// MULTI-STAGE DETECTION
// =====================

// DetectMultiStage checks if attack pattern looks coordinated
func DetectMultiStage(ip string) bool {

	mu.Lock()
	defer mu.Unlock()

	record, exists := attackCache[ip]
	if !exists {
		return false
	}

	// ❌ Already triggered → don't repeat
	if record.triggered {
		return false
	}

	foundSQL := false
	foundMultiple := len(record.events) >= 3

	for _, e := range record.events {
		if e == "SQL_INJECTION" {
			foundSQL = true
		}
	}

	// ✅ Condition: multiple events + SQL attack
	if foundSQL && foundMultiple {
		record.triggered = true // mark triggered
		return true
	}

	return false
}