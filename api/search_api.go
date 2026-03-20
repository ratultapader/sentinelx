package api

import (
	"encoding/json"
	"net/http"

	"sentinelx/search"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	ip := r.URL.Query().Get("ip")
	eventType := r.URL.Query().Get("event")

	results, err := search.SearchEvents(ip, eventType)
	if err != nil {
		http.Error(w, "search failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}