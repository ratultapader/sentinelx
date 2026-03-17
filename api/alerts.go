package api

import (
	"encoding/json"
	"net/http"

	"sentinelx/storage"
)

func AlertsHandler(w http.ResponseWriter, r *http.Request) {

	severity := r.URL.Query().Get("severity")

	alerts, err := storage.GetAlerts(severity)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alerts)
}