package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"sentinelx/storage"
	"sentinelx/ui"
)

type GraphHandler struct {
	graphStore *storage.Neo4jGraphStore
}

func NewGraphHandler(graphStore *storage.Neo4jGraphStore) *GraphHandler {
	return &GraphHandler{
		graphStore: graphStore,
	}
}

func (h *GraphHandler) GetGraphBySourceIP(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	// 🔥 FIX: GET IP FROM PATH (NOT QUERY)
	sourceIP := strings.TrimPrefix(r.URL.Path, "/api/graph/")
	sourceIP = strings.TrimSpace(sourceIP)

	if sourceIP == "" {
		http.Error(w, "missing source_ip", http.StatusBadRequest)
		return
	}

	// ✅ TENANT ENFORCEMENT
	tenantID := strings.TrimSpace(r.Header.Get("X-Tenant-ID"))
	if tenantID == "" {
		http.Error(w, "missing tenant id", http.StatusBadRequest)
		return
	}

	graph, err := h.graphStore.ExportGraphBySourceIP(ctx, sourceIP, tenantID)
	if err != nil {
		http.Error(w, "failed to fetch graph", http.StatusInternalServerError)
		return
	}

	dto := ui.SecurityGraphDTO{
		Nodes: make([]ui.GraphNodeDTO, 0, len(graph.Nodes)),
		Links: make([]ui.GraphLinkDTO, 0, len(graph.Links)),
	}

	for _, n := range graph.Nodes {
		dto.Nodes = append(dto.Nodes, ui.GraphNodeDTO{
			ID:         n.Key,
			Label:      n.Label,
			Name:       graphNodeName(n),
			Properties: n.Properties,
		})
	}

	for _, l := range graph.Links {
		dto.Links = append(dto.Links, ui.GraphLinkDTO{
			Source: l.Source,
			Target: l.Target,
			Type:   l.Type,
			Props:  l.Properties,
		})
	}

	if len(dto.Nodes) == 0 {
		http.Error(w, "graph not found for source_ip", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto)
}

// ================= HELPERS =================

func graphNodeName(n storage.GraphNodeView) string {
	if v, ok := n.Properties["name"].(string); ok && v != "" {
		return v
	}
	if v, ok := n.Properties["ip"].(string); ok && v != "" {
		return v
	}
	if v, ok := n.Properties["endpoint"].(string); ok && v != "" {
		return v
	}
	if v, ok := n.Properties["alert_id"].(string); ok && v != "" {
		return v
	}
	if v, ok := n.Properties["action_type"].(string); ok && v != "" {
		return v
	}
	return n.Key
}