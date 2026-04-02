package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	// "sentinelx/storage"
	"sentinelx/detection"
)

//////////////////////////////////////////////////////
// 🔥 ALERT FETCH (EXISTING)
//////////////////////////////////////////////////////

func AlertsHandler(w http.ResponseWriter, r *http.Request) {
	// severity := r.URL.Query().Get("severity")

	alerts := detection.GetRecentAlerts()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alerts)
}

//////////////////////////////////////////////////////
// 🔥 ALERT WORKFLOW STORE (IN-MEMORY)
//////////////////////////////////////////////////////

type AlertState struct {
	Status string   `json:"status"`
	Notes  []string `json:"notes"`
}

type UpdateRequest struct {
	Status string `json:"status"`
	Note   string `json:"note"`
}

// thread-safe store
var alertStore = struct {
	sync.RWMutex
	data map[string]*AlertState
}{
	data: make(map[string]*AlertState),
}

//////////////////////////////////////////////////////
// 🔥 UPDATE ALERT (NEW API)
//////////////////////////////////////////////////////

func UpdateAlert(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 👉 extract alert ID
	id := strings.TrimPrefix(r.URL.Path, "/api/alerts/")

	if id == "" {
		http.Error(w, "missing alert id", http.StatusBadRequest)
		return
	}

	var req UpdateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	alertStore.Lock()
	defer alertStore.Unlock()

	// 👉 get or create alert state
	state, exists := alertStore.data[id]
	if !exists {
		state = &AlertState{
			Status: "NEW",
			Notes:  []string{},
		}
		alertStore.data[id] = state
	}

	// 👉 update status
	if req.Status != "" {
		state.Status = req.Status
	}

	// 👉 add note
	if req.Note != "" {
		state.Notes = append(state.Notes, req.Note)
	}

	// 👉 response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":     id,
		"status": state.Status,
		"notes":  state.Notes,
	})
}