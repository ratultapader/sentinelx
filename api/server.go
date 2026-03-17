package api

import (
	"fmt"
	"net/http"
	"sentinelx/stream"
)

func StartAPIServer() {

	mux := http.NewServeMux()

	mux.HandleFunc("/metrics", MetricsHandler)
	mux.HandleFunc("/top_attackers", TopAttackersHandler)
	mux.HandleFunc("/health", HealthHandler)

	// ✅ FIXED
	mux.HandleFunc("/dashboard/overview", DashboardOverviewHandler)
	mux.HandleFunc("/dashboard/alerts", DashboardAlertsHandler)
	mux.HandleFunc("/ws", stream.HandleConnections)

	mux.HandleFunc("/events", EventsHandler)
mux.HandleFunc("/alerts", AlertsHandler)

	fmt.Println("SentinelX API running on :9090")

	http.ListenAndServe(":9090", mux)
}