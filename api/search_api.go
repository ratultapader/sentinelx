package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"sentinelx/storage"
)

type SearchHandler struct {
	esStore *storage.ElasticsearchStore
}

func NewSearchHandler(es *storage.ElasticsearchStore) *SearchHandler {
	return &SearchHandler{
		esStore: es,
	}
}

func (h *SearchHandler) SearchHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ip := strings.TrimSpace(r.URL.Query().Get("ip"))
	eventType := strings.TrimSpace(r.URL.Query().Get("event"))

	// ✅ TENANT ENFORCEMENT
	tenantID := strings.TrimSpace(r.Header.Get("X-Tenant-ID"))
	if tenantID == "" {
		http.Error(w, "missing tenant id", http.StatusBadRequest)
		return
	}

	// ✅ FETCH DATA SAFELY
	results := []map[string]interface{}{}

	// If IP provided → use tenant-safe search
	if ip != "" {
		data, err := h.esStore.SearchBySourceIPAndTenant(ctx, storage.IndexSecurityEvents, ip, tenantID)
		if err != nil {
			http.Error(w, "search failed", http.StatusInternalServerError)
			return
		}
		results = data
	} else {
		// fallback → fetch limited tenant data
		data, err := h.esStore.SearchAllByTenant(ctx, storage.IndexSecurityEvents, tenantID, 200)
		if err != nil {
			http.Error(w, "search failed", http.StatusInternalServerError)
			return
		}
		results = data
	}

	// Optional filter by event type
	if eventType != "" {
		filtered := []map[string]interface{}{}
		for _, r := range results {
			if getString(r, "event_type") == eventType {
				filtered = append(filtered, r)
			}
		}
		results = filtered
	}

	// response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}