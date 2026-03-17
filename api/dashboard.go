package api

import (
	"encoding/json"
	"net/http"

	"sentinelx/metrics"
	"sentinelx/detection"
)

func DashboardOverviewHandler(w http.ResponseWriter, r *http.Request) {

	response := map[string]interface{}{
		"total_events":     metrics.TotalEvents,
		"total_alerts":     metrics.TotalAlerts,
		"top_attack_types": metrics.AttackTypes,
		"top_attackers":    metrics.AttackerIPs,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DashboardAlertsHandler(w http.ResponseWriter, r *http.Request) {

	alerts := detection.GetRecentAlerts()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alerts)
}