package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"sentinelx/storage"
)

// 🚀 PRODUCTION-LEVEL ATTACK PATH
func GetAttackPath(w http.ResponseWriter, r *http.Request) {

	// ==========================
	// TENANT
	// ==========================
	tenantID := strings.TrimSpace(r.Header.Get("X-Tenant-ID"))
	if tenantID == "" {
		http.Error(w, "missing tenant id", http.StatusBadRequest)
		return
	}

	// ==========================
	// EXTRACT IP
	// ==========================
	ip := strings.TrimPrefix(r.URL.Path, "/api/attack_path/")
	if ip == "" {
		http.Error(w, "missing ip", http.StatusBadRequest)
		return
	}

	// ==========================
	// FETCH ALERTS (REAL DATA)
	// ==========================
	alerts, err := storage.GetAlerts(tenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// fallback (important for dev)
	if len(alerts) == 0 {
		alerts, _ = storage.GetAlerts("")
	}

	// ==========================
	// BUILD MULTI-STAGE PATH
	// ==========================
	path := []string{ip}

	for _, a := range alerts {
		if a.SourceIP != ip {
			continue
		}

		// 🔹 Attack stage (target)
		if a.Target != "" {
			path = append(path, a.Target)
		}

		// 🔹 Alert stage (with severity)
		if a.Severity != "" {
			path = append(path, "alert:"+a.Severity)
		} else {
			path = append(path, "alert:unknown")
		}

		// 🔹 Response stage (derived from severity)
		switch a.Severity {
		case "critical":
			path = append(path, "response:block_ip")
		case "high":
			path = append(path, "response:rate_limit")
		case "medium":
			path = append(path, "response:monitor")
		default:
			path = append(path, "response:alert_only")
		}
	}

	// ==========================
	// EMPTY CASE
	// ==========================
	if len(path) == 1 {
		path = append(path, "no_activity")
	}

	// ==========================
	// RESPONSE
	// ==========================
	resp := map[string]interface{}{
		"path": path,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}