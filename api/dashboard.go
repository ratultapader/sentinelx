package api

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	"sentinelx/dashboard"
	"sentinelx/detection"
	"sentinelx/storage"
)

//
// ===============================
// 🔥 MAIN DASHBOARD (PRODUCTION)
// ===============================
//

func DashboardOverviewHandler(w http.ResponseWriter, r *http.Request) {

	tenantID := strings.TrimSpace(r.Header.Get("X-Tenant-ID"))
	if tenantID == "" {
		http.Error(w, "missing tenant id", http.StatusBadRequest)
		return
	}

	// ==========================
	// GET ALERTS
	// ==========================
	alerts, err := storage.GetAlerts(tenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// fallback
	if len(alerts) == 0 {
		alerts, _ = storage.GetAlerts("")
	}

	totalAlerts := len(alerts)

	// ==========================
	// AGGREGATIONS
	// ==========================
	attackers := map[string]int{}
	attackTypes := map[string]int{}

	for _, a := range alerts {
		ip := a.SourceIP

		if ip != "" && ip != "::1" && ip != "127.0.0.1" {
			attackers[ip]++
		}

		if a.Target != "" {
			attackTypes[a.Target]++
		}
	}

	// ==========================
	// SORT + LIMIT (TOP 5)
	// ==========================
	type kv struct {
		Key   string
		Value int
	}

	// ---- attackers ----
	var attackerList []kv
	for k, v := range attackers {
		attackerList = append(attackerList, kv{k, v})
	}

	sort.Slice(attackerList, func(i, j int) bool {
		return attackerList[i].Value > attackerList[j].Value
	})

	topAttackers := map[string]int{}
	for i, item := range attackerList {
		if i >= 5 {
			break
		}
		topAttackers[item.Key] = item.Value
	}

	// ---- attack types ----
	var typeList []kv
	for k, v := range attackTypes {
		typeList = append(typeList, kv{k, v})
	}

	sort.Slice(typeList, func(i, j int) bool {
		return typeList[i].Value > typeList[j].Value
	})

	topTypes := map[string]int{}
	for i, item := range typeList {
		if i >= 5 {
			break
		}
		topTypes[item.Key] = item.Value
	}

	// ==========================
	// RESPONSE
	// ==========================
	response := map[string]interface{}{
		"total_events":     totalAlerts,
		"total_alerts":     totalAlerts,
		"top_attack_types": topTypes,
		"top_attackers":    topAttackers,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

//
// ===============================
// 🔥 ALERTS LIST (USED BY UI)
// ===============================
//

func DashboardAlertsHandler(w http.ResponseWriter, r *http.Request) {

	alerts := detection.GetRecentAlerts()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alerts)
}

//
// ===============================
// 🔥 ANOMALY DASHBOARD (ADVANCED)
// ===============================
//

type AnomalyDashboardHandler struct {
	service *dashboard.DashboardService
	store   *storage.ElasticsearchStore
}

func NewAnomalyDashboardHandler(service *dashboard.DashboardService, store *storage.ElasticsearchStore) *AnomalyDashboardHandler {
	return &AnomalyDashboardHandler{
		service: service,
		store:   store,
	}
}

func (h *AnomalyDashboardHandler) GetAnomalyDashboard(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	tenantID := strings.TrimSpace(r.Header.Get("X-Tenant-ID"))
	if tenantID == "" {
		http.Error(w, "missing tenant id", http.StatusBadRequest)
		return
	}

	alerts, err := h.store.SearchAllByTenant(ctx, storage.IndexAlerts, tenantID, 1000)
	if err != nil {
		http.Error(w, "failed to fetch alerts", http.StatusInternalServerError)
		return
	}

	resp, err := h.service.BuildDashboard(ctx, alerts)
	if err != nil {
		http.Error(w, "failed to build dashboard", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}