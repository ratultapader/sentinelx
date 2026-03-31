package collector

import (
	"encoding/json"
	"fmt"
	"net/http"

	"sentinelx/models"
	"sentinelx/multi_tenant"
	"sentinelx/pipeline"
)

func StartHTTPServer() {

	mux := http.NewServeMux()

	// EXISTING
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Login page")
	})

	// EVENT INGESTION
	mux.HandleFunc("/event", func(w http.ResponseWriter, r *http.Request) {

		// TENANT REQUIRED
		tenantID := multi_tenant.TenantIDFromRequest(r)
		if tenantID == "" {
			http.Error(w, "missing tenant id", http.StatusBadRequest)
			return
		}

		var req struct {
			SourceIP  string `json:"source_ip"`
			EventType string `json:"event_type"`
			Payload   string `json:"payload"`
			Target    string `json:"target"`
		}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		// Create proper event
		event := models.NewSecurityEvent(req.EventType)
		event.TenantID = tenantID
		event.SourceIP = req.SourceIP
		event.Protocol = r.Proto

		// SecurityEvent has no Target field in this repo, so keep it in metadata.
		if req.Target != "" {
			event.Metadata["target"] = req.Target
		}

		// CRITICAL LINE
		event.Metadata["payload"] = req.Payload

		ctx := multi_tenant.WithTenantID(r.Context(), tenantID)

		fmt.Println("Sending event:", event.EventType, event.SourceIP)
		pipeline.PublishEvent(ctx, event)
		fmt.Println("Event received:", event.EventType, event.SourceIP)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"event received"}`))
	})

	handler := mux

	fmt.Println("HTTP Server running on :8080")

	http.ListenAndServe(":8080", handler)
}
