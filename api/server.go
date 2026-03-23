package api

import (
	"fmt"
	"net/http"

	"sentinelx/storage"
	"sentinelx/stream"
)

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func StartAPIServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/metrics", MetricsHandler)
	mux.HandleFunc("/top_attackers", TopAttackersHandler)
	mux.HandleFunc("/health", HealthHandler)

	mux.HandleFunc("/dashboard/overview", DashboardOverviewHandler)
	mux.HandleFunc("/dashboard/alerts", DashboardAlertsHandler)
	mux.HandleFunc("/ws", stream.HandleConnections)

	mux.HandleFunc("/events", EventsHandler)
	mux.HandleFunc("/incidents", IncidentsHandler)
	mux.HandleFunc("/alerts", AlertsHandler)
	mux.HandleFunc("/search", SearchHandler)

	if storage.GraphStore != nil {
		graphHandler := NewGraphHandler(storage.GraphStore)
		mux.HandleFunc("/api/graph", graphHandler.GetGraphBySourceIP)
	}

	fmt.Println("SentinelX API running on :9090")
	err := http.ListenAndServe(":9090", withCORS(mux))
	if err != nil {
		fmt.Println("API server error:", err)
	}
}