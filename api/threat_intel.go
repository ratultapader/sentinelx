package api

import (
	"encoding/json"
	"net/http"

	"sentinelx/threatintel"
)

func ThreatIntelHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {

	case "GET":
		data := threatintel.GetAllThreats()
		json.NewEncoder(w).Encode(data)

	case "POST":
		var req struct {
			IP string `json:"ip"`
		}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil || req.IP == "" {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		threatintel.AddThreat(req.IP)

		w.WriteHeader(http.StatusCreated)
	}
}