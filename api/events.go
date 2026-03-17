package api

import (
	"encoding/json"
	"net/http"

	"sentinelx/storage"
)

func EventsHandler(w http.ResponseWriter, r *http.Request) {

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