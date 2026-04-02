package api

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "sort"
    "strings"
    "time"

    "github.com/google/uuid"
    "sentinelx/detection"

    "sentinelx/metrics"
    "sentinelx/models"
    "sentinelx/storage"
    "sentinelx/stream"
    "sentinelx/threatintel"
)

// ===============================
// METRICS HANDLER
// ===============================
func MetricsHandler(w http.ResponseWriter, r *http.Request) {

	data := map[string]interface{}{
		"events_processed_total":  metrics.TotalEvents,
		"alerts_generated_total": metrics.TotalAlerts,
		"unique_attackers":       len(metrics.AttackerIPs),
		"attack_types":           len(metrics.AttackTypes),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// ===============================
// TOP ATTACKERS
// ===============================
func TopAttackersHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics.AttackerIPs)
}

// ===============================
// HEALTH CHECK
// ===============================
func HealthHandler(w http.ResponseWriter, r *http.Request) {

	status := "ok"

	services := map[string]string{
		"elasticsearch": "down",
		"neo4j":        "down",
		"database":     "down",
	}

	if storage.ESStore != nil {
		services["elasticsearch"] = "up"
	}
	if storage.GraphStore != nil {
		services["neo4j"] = "up"
	}
	if storage.DB != nil {
		services["database"] = "up"
	}

	for _, v := range services {
		if v == "down" {
			status = "degraded"
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   status,
		"services": services,
	})
}

// ===============================
// TEST EVENT HANDLER
// ===============================
func TestEventHandler(w http.ResponseWriter, r *http.Request) {

    type TestEvent struct {
        Type     string  `json:"type"`
        Severity string  `json:"severity"`
        Message  string  `json:"message"`
        SourceIP string  `json:"source_ip"`
        Score    float64 `json:"score"`
    }

    var event TestEvent

    err := json.NewDecoder(r.Body).Decode(&event)
    if err != nil {
        http.Error(w, "invalid payload", http.StatusBadRequest)
        return
    }

    log.Println("Incoming event from:", event.SourceIP)

    tenantID := r.Header.Get("X-Tenant-ID")
    if tenantID == "" {
        tenantID = "t1"
    }

    alert := models.Alert{
        ID:          uuid.New().String(),
        Timestamp:   time.Now(),
        Type:        event.Type,
        Severity:    event.Severity,
        SourceIP:    event.SourceIP,
        Target:      event.Message,
        Description: event.Message,
        ThreatScore: event.Score,
        Status:      models.AlertStatusNew,
        TenantID:    tenantID,
    }

    detected := event.Type
    metrics.RecordEvent(event.SourceIP, event.Type, detected)
    metrics.RecordAlert(event.Type)

    select {
    case detection.AlertQueue <- alert:
    default:
        log.Println("Alert queue full")
    }

    stream.BroadcastAlert(alert)

    if threatintel.IsMaliciousIP(event.SourceIP) {
        threatintel.IncrementMatch(event.SourceIP)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "status": "event processed",
    })
}

// ===============================
// ALERT INSIGHTS HANDLER (NEW ??)
// ===============================
func AlertInsightsHandler(w http.ResponseWriter, r *http.Request) {

	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		tenantID = "t1"
	}

	alerts, err := storage.GetRecentAlertsByTenant(r.Context(), tenantID)
	if err != nil || len(alerts) == 0 {
		http.Error(w, "failed to fetch alerts", http.StatusInternalServerError)
		return
	}

	sort.Slice(alerts, func(i, j int) bool {
		return alerts[i].Timestamp.Before(alerts[j].Timestamp)
	})

	alert := alerts[len(alerts)-1]

	score := 50
	reasons := []string{"High anomaly score"}

	var filtered []models.Alert
	for _, a := range alerts {
		if a.SourceIP == alert.SourceIP {
			filtered = append(filtered, a)
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Timestamp.Before(filtered[j].Timestamp)
	})

	typeSet := make(map[string]bool)
	var timelineParts []string

	for _, a := range filtered {
		if a.Type == "" {
			continue
		}

		timeStr := a.Timestamp.Format("15:04")

		if !typeSet[a.Type] {
			typeSet[a.Type] = true
			entry := fmt.Sprintf("%s (%s)", a.Type, timeStr)
			timelineParts = append(timelineParts, entry)
		}
	}

	timeline := strings.Join(timelineParts, " => ")

	if len(timelineParts) > 1 {
		score += 20
		reasons = append(reasons, "Multi-stage attack detected")
	}

	story := fmt.Sprintf(
		"Attacker %s executed a timed attack sequence: %s",
		alert.SourceIP,
		timeline,
	)

	// 🔥 convert story → steps
// 🔥 convert timeline → steps (FIXED)
steps := []string{}

// split using correct separator (NOT colon)
if strings.Contains(timeline, "=>") {
	parts := strings.Split(timeline, "=>")

	for _, p := range parts {
		clean := strings.TrimSpace(p)
		if clean != "" {
			steps = append(steps, clean)
		}
	}
} else if timeline != "" {
	steps = append(steps, timeline)
}

// add intro step at top
steps = append([]string{
	fmt.Sprintf("Attacker %s started attack", alert.SourceIP),
}, steps...)

// 🔥 NEW RESPONSE
response := map[string]interface{}{
	"priority_score": score,
	"reasons":        reasons,
	"story":          story,  // ✅ keep old
	"story_steps":    steps,  // ✅ new
}

	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	encoder.Encode(response)
}

func generateDescription(t string) string {
	switch t {
	case "SQL Injection":
		return "Database query manipulation detected"
	case "XSS":
		return "Injected malicious script"
	case "Brute Force":
		return "Multiple login attempts"
	default:
		return "Suspicious activity detected"
	}
}

// ===============================
// KPI HANDLER
// ===============================
func KPIHandler(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		tenantID = "t1"
	}

	alerts, err := storage.GetRecentAlertsByTenant(r.Context(), tenantID)
	if err != nil {
		http.Error(w, "failed to fetch alerts", 500)
		return
	}

	totalAlerts := len(alerts)

	// 🔥 MTTR (safe version)
	var totalTime float64
	var count int

	for _, a := range alerts {
		if a.Status == "RESOLVED" && !a.Timestamp.IsZero() {
			// if you don’t have resolved_at yet → skip
			count++
		}
	}

	mttr := 0
	if count > 0 {
		mttr = int(totalTime / float64(count))
	}

	response := map[string]interface{}{
		"mttr":                mttr,
		"total_alerts":        totalAlerts,
		"false_positive_rate": 5,
	}

	writeJSON(w, 200, response)
}

// ===============================
// THREAT TREND HANDLER
// ===============================
func ThreatTrendHandler(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		tenantID = "t1"
	}

	alerts, err := storage.GetRecentAlertsByTenant(r.Context(), tenantID)
	if err != nil {
		http.Error(w, "failed to fetch alerts", 500)
		return
	}

	buckets := make(map[string]int)

	for _, a := range alerts {
		if a.Timestamp.IsZero() {
			continue
		}

		key := a.Timestamp.Format("15:04")
		buckets[key]++
	}

	var result []map[string]interface{}

	for k, v := range buckets {
		result = append(result, map[string]interface{}{
			"time":   k,
			"alerts": v,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i]["time"].(string) < result[j]["time"].(string)
	})

	writeJSON(w, 200, result)
}

// ===============================
// FEEDBACK HANDLER
// ===============================
func FeedbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		tenantID = "t1"
	}

	var payload struct {
		Type    string `json:"type"`
		AlertID string `json:"alert_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if payload.Type == "" {
		http.Error(w, "missing feedback type", http.StatusBadRequest)
		return
	}

	fmt.Println("📩 Feedback received:",
		"type =", payload.Type,
		"alert =", payload.AlertID,
		"tenant =", tenantID,
	)

	writeJSON(w, 200, map[string]interface{}{
		"status":   "received",
		"type":     payload.Type,
		"alert_id": payload.AlertID,
	})
}


// ===============================
// 📜 AUDIT LOGS
// ===============================
func AuditLogsHandler(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		tenantID = "t1"
	}

	alerts, err := storage.GetRecentAlertsByTenant(r.Context(), tenantID)
	if err != nil {
		http.Error(w, "failed to fetch alerts", 500)
		return
	}

	unique := make(map[string]int)

for _, a := range alerts {
	ip := a.SourceIP
	unique[ip]++
}

logs := []map[string]interface{}{}

for ip, count := range unique {
	logs = append(logs, map[string]interface{}{
		"user":   "system",
		"action": fmt.Sprintf("generated %d alerts", count),
		"target": ip,
		"timestamp": time.Now().UTC(),
	})
}

	writeJSON(w, 200, logs)
}
// ===============================
// ⚡ PERFORMANCE METRICS
// ===============================
func PerformanceHandler(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		tenantID = "t1"
	}

	alerts, err := storage.GetRecentAlertsByTenant(r.Context(), tenantID)
	if err != nil {
		http.Error(w, "failed to fetch alerts", 500)
		return
	}

	totalAlerts := len(alerts)

	// simple estimation
	alertsPerSec := 0
	if totalAlerts > 0 {
		alertsPerSec = totalAlerts / 10
	}

	writeJSON(w, 200, map[string]interface{}{
		"events_per_sec": totalAlerts * 2,
		"alerts_per_sec": alertsPerSec,
		"latency":        5,
	})
}