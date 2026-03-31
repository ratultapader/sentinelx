package api

import (
	"fmt"
	"net/http"
	"strings"

	"sentinelx/dashboard"
	"sentinelx/storage"
	"sentinelx/stream"
)

// ===============================
// START API SERVER
// ===============================
func StartAPIServer() {
	mux := http.NewServeMux()

	// ===============================
// GRAPH HANDLER INIT
// ===============================
if storage.GraphStore != nil {
	graphHandler := NewGraphHandler(storage.GraphStore)
	mux.HandleFunc("/api/graph/", graphHandler.GetGraphBySourceIP)
}

	// ===============================
	// CORE APIs
	// ===============================
	mux.HandleFunc("/metrics", MetricsHandler)
	mux.HandleFunc("/top_attackers", TopAttackersHandler)
	mux.HandleFunc("/health", HealthHandler)

	mux.HandleFunc("/api/alerts/recent", GetRecentAlerts)
	mux.HandleFunc("/api/alerts", GetRecentAlerts)
	mux.HandleFunc("/api/alerts/", UpdateAlert)

	mux.HandleFunc("/api/threat_intel", ThreatIntelHandler)
	mux.HandleFunc("/api/alert_insights", AlertInsightsHandler)
	mux.HandleFunc("/api/attack_map", GetAttackMap)
	mux.HandleFunc("/api/attack_path/", GetAttackPath)
	mux.HandleFunc("/api/attack_pattern/", GetAttackPattern)

	// ===============================
	// RULES
	// ===============================
	mux.HandleFunc("/api/rules", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			GetRules(w, r)
		} else if r.Method == "POST" {
			CreateRule(w, r)
		}
	})
	mux.HandleFunc("/api/rules/", ToggleRule)

	// ===============================
	// PLAYBOOK APIs (FIXED)
	// ===============================
	mux.HandleFunc("/api/playbooks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			GetPlaybooks(w, r)
		case "POST":
			CreatePlaybook(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// HANDLE TOGGLE + DELETE
	mux.HandleFunc("/api/playbooks/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			TogglePlaybook(w, r)
		case "DELETE":
			DeletePlaybook(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// ===============================
	// DASHBOARD
	// ===============================
	mux.HandleFunc("/dashboard/overview", DashboardOverviewHandler)
	mux.HandleFunc("/dashboard/alerts", DashboardAlertsHandler)

	// ===============================
	// ANOMALY DASHBOARD
	// ===============================
	if storage.ESStore != nil {
		dashboardService := dashboard.NewDashboardService()
		dashboardHandler := NewAnomalyDashboardHandler(dashboardService, storage.ESStore)
		mux.HandleFunc("/api/dashboard/anomalies", dashboardHandler.GetAnomalyDashboard)
	}

	// ===============================
	// STREAM
	// ===============================
	mux.HandleFunc("/ws", stream.HandleConnections)

	// ===============================
	// INCIDENT APIs
	// ===============================
	if storage.ESStore != nil && storage.GraphStore != nil {
		incidentHandler := NewIncidentHandler(storage.ESStore, storage.GraphStore)

		mux.HandleFunc("/api/incidents", incidentHandler.GetIncidents)
		mux.HandleFunc("/api/incidents/", incidentHandler.GetIncidentByID)
		mux.HandleFunc("/api/timeline/", incidentHandler.GetTimelineByIP)
		// mux.HandleFunc("/api/graph/", incidentHandler.GetGraphByIP)
	}

	// ===============================
	// REPORT APIs
	// ===============================
	if storage.ESStore != nil {
		reportHandler := NewReportHandler(storage.ESStore)

		mux.HandleFunc("/api/reports", reportHandler.ListReports)

		mux.HandleFunc("/api/reports/", func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.URL.Path, "/pdf") {
				reportHandler.GetReportPDF(w, r)
			} else {
				reportHandler.GetReportJSON(w, r)
			}
		})
	}

	// ===============================
	// FINAL MIDDLEWARE (FIXED)
	// ===============================
	handler := Chain(
		mux,
		RecoverMiddleware,
		LoggingMiddleware,
		CORSMiddleware,
		TenantMiddleware,
	)

	fmt.Println("SentinelX API running on :9090")

	err := http.ListenAndServe(":9090", handler)
	if err != nil {
		fmt.Println("API server error:", err)
	}
}
