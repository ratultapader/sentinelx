package api

import (
	"encoding/json"
	"net/http"

	"sentinelx/incident"
)

func IncidentsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(incident.GetAllIncidents())
}