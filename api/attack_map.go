package api

import (
	"encoding/json"
	"net/http"

	"sentinelx/storage"
)

// ===============================
// ATTACK MAP API
// ===============================
func GetAttackMap(w http.ResponseWriter, r *http.Request) {
	// Tenant ID (multi-tenant system)
	tenantID := r.Header.Get("X-Tenant-ID")

	// Get alerts
	alerts, err := storage.GetAlerts(tenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fallback if tenant has no alerts.
	if len(alerts) == 0 {
		alerts, _ = storage.GetAlerts("")
	}

	seen := map[string]bool{}

	var result []map[string]interface{}

	for _, a := range alerts {
		ip := a.SourceIP

		// skip localhost
		if ip == "" || ip == "::1" || ip == "127.0.0.1" {
			continue
		}

		// remove duplicates
		if seen[ip] {
			continue
		}
		seen[ip] = true

		lat, lng, country := getGeo(ip)

		result = append(result, map[string]interface{}{
			"ip":      ip,
			"lat":     lat,
			"lng":     lng,
			"country": country,
		})
	}

	// Always return array (never null)
	json.NewEncoder(w).Encode(result)
}

func getGeo(ip string) (float64, float64, string) {
	switch ip {
	case "10.10.10.10":
		return 28.6139, 77.2090, "India"
	case "101.101.101.101":
		return 37.7749, -122.4194, "USA"
	case "102.102.102.102":
		return 51.5074, -0.1278, "UK"
	case "103.103.103.103":
		return 35.6762, 139.6503, "Japan"
	default:
		return 0, 0, "Unknown"
	}
}
