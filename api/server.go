package api

import (
	"fmt"
	"net/http"
)

func StartAPIServer() {

	mux := http.NewServeMux()

	mux.HandleFunc("/metrics", MetricsHandler)
	mux.HandleFunc("/top_attackers", TopAttackersHandler)
	mux.HandleFunc("/health", HealthHandler)

	fmt.Println("SentinelX API running on :9090")

	http.ListenAndServe(":9090", mux)
}