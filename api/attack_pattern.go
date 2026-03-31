package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"sentinelx/storage"
)

type AttackStep struct {
	Step      string `json:"step"`
	Timestamp string `json:"timestamp"`
}

func GetAttackPattern(w http.ResponseWriter, r *http.Request) {

	tenantID := strings.TrimSpace(r.Header.Get("X-Tenant-ID"))
	if tenantID == "" {
		http.Error(w, "missing tenant id", http.StatusBadRequest)
		return
	}

	ip := strings.TrimPrefix(r.URL.Path, "/api/attack_pattern/")
	if ip == "" {
		http.Error(w, "missing ip", http.StatusBadRequest)
		return
	}

	alerts, err := storage.GetAlerts(tenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// fallback
	if len(alerts) == 0 {
		alerts, _ = storage.GetAlerts("")
	}

	var result []AttackStep

	for _, a := range alerts {
		if a.SourceIP == ip {
			result = append(result, AttackStep{
				Step:      a.Target,
				Timestamp: a.Timestamp,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}