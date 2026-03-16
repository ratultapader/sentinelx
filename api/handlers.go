package api

import (
	"encoding/json"
	"net/http"

	"sentinelx/metrics"
)

func MetricsHandler(w http.ResponseWriter, r *http.Request) {

	data := map[string]interface{}{
		"total_events": metrics.TotalEvents,
		"total_alerts": metrics.TotalAlerts,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func TopAttackersHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics.AttackerIPs)
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {

	data := map[string]string{
		"status": "ok",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}