package api

import (
	// "encoding/json"
	"fmt"
	"net/http"
	"strings"

	"sentinelx/investigation"
	"sentinelx/reporting"
	"sentinelx/storage"
)

type ReportHandler struct {
	es *storage.ElasticsearchStore
}

func NewReportHandler(es *storage.ElasticsearchStore) *ReportHandler {
	return &ReportHandler{es: es}
}

// ================= GET /reports/:id =================

func (h *ReportHandler) GetReportJSON(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	path := r.URL.Path
	path = strings.TrimPrefix(path, "/api/reports/")
	path = strings.TrimPrefix(path, "/reports/")
	id := path

	// Tenant enforcement
	tenantID := strings.TrimSpace(r.Header.Get("X-Tenant-ID"))
	if tenantID == "" {
		tenantID = strings.TrimSpace(r.URL.Query().Get("tenant_id"))
	}
	if tenantID == "" {
		writeError(w, 400, "missing tenant id")
		return
	}

	fmt.Println("DEBUG >>> REPORT ID:", id)
	fmt.Println("DEBUG >>> TENANT:", tenantID)
	fmt.Println("DEBUG >>> CLEAN ID:", id)

	alert, err := h.es.GetByDocumentIDAndTenant(ctx, storage.IndexAlerts, id, tenantID)
	if err != nil || alert == nil {
		fmt.Println("DEBUG >>> ALERT NOT FOUND IN ES ?")
		writeError(w, 404, "incident not found")
		return
	}

	fmt.Println("DEBUG >>> ALERT FOUND ?", alert)

	sourceIP := getString(alert, "source_ip")

	alerts, _ := h.es.SearchBySourceIPAndTenant(ctx, storage.IndexAlerts, sourceIP, tenantID)
	fmt.Println("DEBUG >>> SOURCE IP:", sourceIP)
	fmt.Println("DEBUG >>> RELATED ALERTS COUNT:", len(alerts))

	events := []investigation.TimelineEvent{}
	for _, doc := range alerts {
		events = append(events, investigation.NormalizeAlert(doc))
	}

	timeline := investigation.NewBuilder().Build(sourceIP, events)

	builder := reporting.NewReportBuilder()
	report := builder.Build(id, alert, timeline, nil, nil)

	writeJSON(w, 200, report)
	fmt.Println("DEBUG FINAL REPORT >>>", report)
}

// ================= GET /reports/:id/pdf =================

func (h *ReportHandler) GetReportPDF(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	path := r.URL.Path
	path = strings.TrimPrefix(path, "/api/reports/")
	path = strings.TrimPrefix(path, "/reports/")
	path = strings.TrimSuffix(path, "/pdf")
	id := path

	// Tenant enforcement
	tenantID := strings.TrimSpace(r.Header.Get("X-Tenant-ID"))
	if tenantID == "" {
		tenantID = strings.TrimSpace(r.URL.Query().Get("tenant_id"))
	}
	if tenantID == "" {
		writeError(w, 400, "missing tenant id")
		return
	}

	alert, err := h.es.GetByDocumentIDAndTenant(ctx, storage.IndexAlerts, id, tenantID)
	if err != nil || alert == nil {
		writeError(w, 404, "incident not found")
		return
	}

	sourceIP := getString(alert, "source_ip")

	alerts, _ := h.es.SearchBySourceIPAndTenant(ctx, storage.IndexAlerts, sourceIP, tenantID)

	events := []investigation.TimelineEvent{}
	for _, doc := range alerts {
		events = append(events, investigation.NormalizeAlert(doc))
	}

	timeline := investigation.NewBuilder().Build(sourceIP, events)
	builder := reporting.NewReportBuilder()
	report := builder.Build(id, alert, timeline, nil, nil)

	filePath := "/tmp/report_" + id + ".pdf"

	err = reporting.ExportIncidentReportPDF(report, filePath)
	if err != nil {
		writeError(w, 500, "failed to generate pdf")
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=incident_report.pdf")

	http.ServeFile(w, r, filePath)
}

// ================= GET /api/reports =================

func (h *ReportHandler) ListReports(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	// Tenant
	tenantID := strings.TrimSpace(r.Header.Get("X-Tenant-ID"))
	if tenantID == "" {
		writeError(w, 400, "missing tenant id")
		return
	}

	alerts, err := h.es.SearchAllByTenant(ctx, storage.IndexAlerts, tenantID, 50)
	if err != nil {
		writeError(w, 500, "failed to fetch reports")
		return
	}

	items := []map[string]interface{}{}

	for _, a := range alerts {
		id := getString(a, "id")
		if id == "" {
			id = getString(a, "_id")
		}

		items = append(items, map[string]interface{}{
			"id":          id,
			"incident_id": id,
			"severity":    getString(a, "severity"),
			"created_at":  getString(a, "timestamp"),
		})
	}

	writeJSON(w, 200, map[string]interface{}{
		"items": items,
	})
}

// func GetReportDetail(w http.ResponseWriter, r *http.Request) {

// 	tenantID := r.Header.Get("X-Tenant-ID")
// 	if tenantID == "" {
// 		http.Error(w, "missing tenant id", http.StatusBadRequest)
// 		return
// 	}

// 	// extract ID from URL
// 	id := strings.TrimPrefix(r.URL.Path, "/api/reports/")

// 	// TEMP (mock data for now)
// 	report := map[string]interface{}{
// 		"incident_id": id,
// 		"attack_chain": []map[string]interface{}{
// 			{
// 				"timestamp":  "2026-03-27 10:00:00",
// 				"event_type": "port_scan",
// 				"summary":    "Multiple ports scanned",
// 			},
// 		},
// 		"mitre_tactic":    "Initial Access",
// 		"mitre_technique": "Suspicious Network Activity",
// 	}

// 	json.NewEncoder(w).Encode(report)
// }
