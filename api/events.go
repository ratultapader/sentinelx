package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"sentinelx/detection"
	"sentinelx/storage"
)

func EventsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DEBUG HANDLER >>> TENANT IN CTX:", r.Context().Value("tenant_id"))

	if r.Method == http.MethodPost {
		var req struct {
			EventType string `json:"event_type"`
			SourceIP  string `json:"source_ip"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		ctx := r.Context()

		detection.GenerateAlert(ctx, req.EventType, req.SourceIP, "Detected "+req.EventType)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"event processed"}`))
		return
	}

	ip := r.URL.Query().Get("ip")
	eventType := r.URL.Query().Get("type")

	events, err := storage.GetEvents(ip, eventType)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}
