package api

import (
	"context"
	"net/http"
	"time"

	"sentinelx/storage"
)

func GetRecentAlerts(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")

	if tenantID == "" {
		WriteError(w, http.StatusBadRequest, "missing tenant id")
		return
	}

	if storage.ESStore == nil {
		WriteError(w, http.StatusInternalServerError, "ES not initialized")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	alerts, err := storage.ESStore.SearchAllByTenant(
		ctx,
		storage.IndexAlerts,
		tenantID,
		10,
	)

	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 🔥 FIX STARTS HERE
	var response []map[string]interface{}

	for _, a := range alerts {

		meta := map[string]interface{}{}
		if m, ok := a["metadata"].(map[string]interface{}); ok {
			meta = m
		}

		response = append(response, map[string]interface{}{
			"source_ip":     a["source_ip"],
			"severity":      a["severity"],
			"threat_score":  a["threat_score"],
			"timestamp":     a["timestamp"],
			"target":        a["target"],

			// 🔥 CRITICAL FIX (flatten metadata)
			"anomaly_score":  meta["anomaly_score"],
			"signature_score": meta["signature_match"],
			"ip_reputation":  meta["ip_reputation"],
			"behavior_score": meta["behavior_deviation"],
		})
	}

	WriteJSON(w, http.StatusOK, response)
}