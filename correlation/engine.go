package correlation

import (
	"time"
	"sync"

	"sentinelx/models"
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

// RecordEvent stores attack activity for correlation.
func RecordEvent(ip string, eventType string) {
	mu.Lock()
	defer mu.Unlock()

	record, exists := attackCache[ip]

	// If new IP -> create record
	if !exists {
		attackCache[ip] = &attackRecord{
			events:    []string{eventType},
			lastSeen:  time.Now(),
			triggered: false,
		}
		return
	}

	// Reset window after 30 seconds
	if time.Since(record.lastSeen) > 30*time.Second {
		record.events = []string{eventType}
		record.triggered = false
		record.lastSeen = time.Now()
		return
	}

	// Append event
	record.events = append(record.events, eventType)
	record.lastSeen = time.Now()

	// Keep only last 5 events
	if len(record.events) > 5 {
		record.events = record.events[1:]
	}
}

// =====================
// MULTI-STAGE DETECTION
// =====================

// DetectMultiStage checks if attack pattern looks coordinated.
func DetectMultiStage(ip string) bool {
	mu.Lock()
	defer mu.Unlock()

	record, exists := attackCache[ip]
	if !exists {
		return false
	}

	// already triggered
	if record.triggered {
		return false
	}

	foundSQL := false
	foundMultiple := len(record.events) >= 3

	for _, e := range record.events {
		if e == "SQL_INJECTION" || e == "sql_injection" {
			foundSQL = true
		}
	}

	// multiple suspicious events + SQL activity
	if foundSQL && foundMultiple {
		record.triggered = true
		return true
	}

	return false
}

// BuildMultiStageAlert creates a standardized multi-stage attack alert.
func BuildMultiStageAlert(ip string, tenantID string) *models.Alert {
	mu.Lock()
	record, exists := attackCache[ip]
	if !exists {
		mu.Unlock()
		return nil
	}

	eventsCopy := make([]string, len(record.events))
	copy(eventsCopy, record.events)
	mu.Unlock()

	alert := models.Alert{
		ID:          generateCorrelationAlertID(),
		TenantID: tenantID,
		Timestamp:   time.Now().UTC(),
		Type:        "multi_stage_attack",
		Severity:    models.SeverityCritical,
		SourceIP:    ip,
		Description: "Possible coordinated multi-stage attack detected",
		ThreatScore: 0.99,
		Status:      models.AlertStatusNew,
		Metadata: map[string]interface{}{
			"correlated_events": eventsCopy,
			"event_count":       len(eventsCopy),
		},
	}

	return &alert
}

// generateCorrelationAlertID generates a unique alert ID for correlation alerts.
func generateCorrelationAlertID() string {
	return "ALT-" + time.Now().UTC().Format("20060102150405.000000000")
}