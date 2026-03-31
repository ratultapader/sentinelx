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

	response := map[string]interface{}{
		"priority_score": score,
		"reasons":        reasons,
		"story":          story,
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





