package api

import (
	// "context"
	"encoding/json"
	"net/http"
	"strings"
	"sort"
	"time"

	"sentinelx/incident"
	"sentinelx/storage"
)

// ================= OLD HANDLER (KEEP FOR NOW) =================

func IncidentsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(incident.GetAllIncidents())
}

// ================= NEW PRODUCTION HANDLER =================

type IncidentHandler struct {
	esStore    *storage.ElasticsearchStore
	graphStore *storage.Neo4jGraphStore
}

func NewIncidentHandler(es *storage.ElasticsearchStore, graph *storage.Neo4jGraphStore) *IncidentHandler {
	return &IncidentHandler{
		esStore:    es,
		graphStore: graph,
	}
}

// ================= GET /incidents =================

// ================= GET /incidents =================

func (h *IncidentHandler) GetIncidents(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	tenantID := strings.TrimSpace(r.Header.Get("X-Tenant-ID"))
	if tenantID == "" {
		writeError(w, 400, "missing tenant id")
		return
	}

	// ===============================
	// 🔥 REDIS CACHE KEY
	// ===============================
	cacheKey := "incidents:" + tenantID

	// ===============================
	// 🔥 STEP 1: CHECK CACHE
	// ===============================
	cached, err := storage.RDB.Get(ctx, cacheKey).Result()
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cached))
		return
	}

	// ===============================
	// 🔥 STEP 2: FETCH FROM ES
	// ===============================
	alerts, err := h.esStore.SearchAllByTenant(ctx, storage.IndexAlerts, tenantID, 500)
	if err != nil {
		writeError(w, 500, "failed to fetch incidents")
		return
	}

	grouped := make(map[string][]map[string]interface{})

	for _, doc := range alerts {
		ip := getString(doc, "source_ip")
		if ip == "" {
			continue
		}
		grouped[ip] = append(grouped[ip], doc)
	}

	items := []map[string]interface{}{}

	for ip, group := range grouped {
		latest := group[0]

		for _, g := range group {
			if getString(g, "timestamp") > getString(latest, "timestamp") {
				latest = g
			}
		}

		severity := "low"

		for _, a := range group {
			s := getString(a, "severity")

			if s == "critical" {
				severity = "critical"
				break
			} else if s == "high" && severity != "critical" {
				severity = "high"
			} else if s == "medium" && severity == "low" {
				severity = "medium"
			}
		}

		items = append(items, map[string]interface{}{
			"id":          ip,
			"source_ip":   ip,
			"timestamp":   getString(latest, "timestamp"),
			"severity":    severity,
			"alert_count": len(group),
			"alerts":      group,
		})
	}

	response := map[string]interface{}{
		"items": items,
		"count": len(items),
	}

	jsonData, _ := json.Marshal(response)

	// ===============================
	// 🔥 STEP 3: STORE IN REDIS
	// ===============================
	storage.RDB.Set(ctx, cacheKey, jsonData, 10*time.Second)

	// ===============================
	// 🔥 RETURN RESPONSE
	// ===============================
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
// ================= GET /incidents/:id =================

func (h *IncidentHandler) GetIncidentByID(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ip := strings.TrimPrefix(r.URL.Path, "/api/incidents/")

	tenantID := strings.TrimSpace(r.Header.Get("X-Tenant-ID"))
	if tenantID == "" {
		writeError(w, 400, "missing tenant id")
		return
	}

	if ip == "" {
		writeError(w, 400, "missing incident id")
		return
	}

	alerts, err := h.esStore.SearchBySourceIPAndTenant(ctx, storage.IndexAlerts, ip, tenantID)
	if err != nil {
		writeError(w, 500, "failed to fetch incident")
		return
	}

	if len(alerts) == 0 {
		writeError(w, 404, "incident not found")
		return
	}

	severity := "low"

	for _, a := range alerts {
		s := getString(a, "severity")

		if s == "critical" {
			severity = "critical"
			break
		} else if s == "high" && severity != "critical" {
			severity = "high"
		} else if s == "medium" && severity == "low" {
			severity = "medium"
		}
	}

	response := map[string]interface{}{
		"id":          ip,
		"source_ip":   ip,
		"severity":    severity,
		"alert_count": len(alerts),
		"alerts":      alerts,
	}

	writeJSON(w, 200, response)
}

// ================= GET /timeline/:ip =================

func (h *IncidentHandler) GetTimelineByIP(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ip := strings.TrimPrefix(r.URL.Path, "/api/timeline/")

	// ✅ TENANT ENFORCEMENT
	tenantID := strings.TrimSpace(r.Header.Get("X-Tenant-ID"))
	if tenantID == "" {
		writeError(w, 400, "missing tenant id")
		return
	}

	alerts, _ := h.esStore.SearchBySourceIPAndTenant(ctx, storage.IndexAlerts, ip, tenantID)
	events, _ := h.esStore.SearchBySourceIPAndTenant(ctx, storage.IndexSecurityEvents, ip, tenantID)
	actions, _ := h.esStore.SearchBySourceIPAndTenant(ctx, storage.IndexResponseActions, ip, tenantID)

	timeline := []map[string]interface{}{}

// 🔹 EVENTS
for _, e := range events {
	timestamp := getString(e, "timestamp")
	eventType := getString(e, "event_type")

	if timestamp != "" {
		timeline = append(timeline, map[string]interface{}{
			"type":      eventType,
			"timestamp": timestamp,
		})
	}
}

// 🔹 ALERTS
for _, a := range alerts {
	timestamp := getString(a, "timestamp")
	alertType := getString(a, "type")

	if timestamp != "" {
		timeline = append(timeline, map[string]interface{}{
			"type":      alertType,
			"timestamp": timestamp,
		})
	}
}

// 🔹 ACTIONS
for _, ac := range actions {
	timestamp := getString(ac, "timestamp")
	actionType := getString(ac, "action_type")

	if timestamp != "" {
		timeline = append(timeline, map[string]interface{}{
			"type":      actionType,
			"timestamp": timestamp,
		})
	}
}

// 🔥 SORT BY TIME (STRING SORT SAFE FOR ISO TIME)
sort.Slice(timeline, func(i, j int) bool {
	return timeline[i]["timestamp"].(string) < timeline[j]["timestamp"].(string)
})

// 🔥 FINAL RESPONSE (ONLY TIMELINE)
writeJSON(w, 200, timeline)
}

// ================= GET /graph/:ip =================

func (h *IncidentHandler) GetGraphByIP(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ip := strings.TrimPrefix(r.URL.Path, "/api/graph/")

	// ? TENANT
	tenantID := strings.TrimSpace(r.Header.Get("X-Tenant-ID"))
	if tenantID == "" {
		writeError(w, 400, "missing tenant id")
		return
	}

	// ?? FETCH ALERTS FROM ELASTICSEARCH
	alerts, err := h.esStore.SearchBySourceIPAndTenant(ctx, storage.IndexAlerts, ip, tenantID)
	if err != nil {
		writeError(w, 500, "failed to fetch alerts")
		return
	}

	nodes := []map[string]interface{}{}
	edges := []map[string]interface{}{}

	// ? ADD MAIN IP NODE
	nodes = append(nodes, map[string]interface{}{
		"id":   ip,
		"type": "ip",
	})

	// ? BUILD GRAPH FROM ALERTS
	for _, a := range alerts {

		dest := getString(a, "destination")

		if dest == "" {
			continue
		}

		// ADD DEST NODE
		nodes = append(nodes, map[string]interface{}{
			"id":   dest,
			"type": "api",
		})

		// ADD EDGE
		edges = append(edges, map[string]interface{}{
			"source": ip,
			"target": dest,
		})
	}

	// ? FINAL RESPONSE
	response := map[string]interface{}{
		"nodes": nodes,
		"edges": edges,
	}

	writeJSON(w, 200, response)
}

// ================= HELPERS =================

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, ErrorResponse{Error: msg})
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func getFloat(m map[string]interface{}, key string) float64 {
	if v, ok := m[key].(float64); ok {
		return v
	}
	return 0
}





