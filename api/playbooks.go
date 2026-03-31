package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	// "strconv"
	"strings"
)

// ===============================
// MODEL
// ===============================
type Playbook struct {
	ID        string `json:"id"`
	Condition string `json:"condition"`
	Action    string `json:"action"`
	Enabled   bool   `json:"enabled"`
}

// ===============================
// IN-MEMORY STORE (OK FOR NOW)
// ===============================
var playbooks = []Playbook{
	{"p1", "threat_score > 0.9", "block_ip", true},
	{"p2", "threat_score > 0.7", "rate_limit", true},
	{"p3", "threat_score > 0.5", "alert", true},
}

// ===============================
// GET ALL
// ===============================
func GetPlaybooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(playbooks); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ===============================
// CREATE
// ===============================
func CreatePlaybook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var p Playbook

	// decode safely
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// validation
	if p.Condition == "" || p.Action == "" {
		http.Error(w, "condition and action required", http.StatusBadRequest)
		return
	}

	// 🔥 FIXED ID GENERATION
	newID := fmt.Sprintf("p%d", len(playbooks)+1)
	p.ID = newID
	p.Enabled = true

	playbooks = append(playbooks, p)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

// ===============================
// TOGGLE
// ===============================
func TogglePlaybook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := strings.TrimPrefix(r.URL.Path, "/api/playbooks/")

	for i := range playbooks {
		if playbooks[i].ID == id {
			playbooks[i].Enabled = !playbooks[i].Enabled
			json.NewEncoder(w).Encode(playbooks[i])
			return
		}
	}

	http.Error(w, "playbook not found", http.StatusNotFound)
}

// ===============================
// DELETE 🔥 (IMPORTANT)
// ===============================
func DeletePlaybook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := strings.TrimPrefix(r.URL.Path, "/api/playbooks/")

	for i := range playbooks {
		if playbooks[i].ID == id {

			// remove item
			playbooks = append(playbooks[:i], playbooks[i+1:]...)

			json.NewEncoder(w).Encode(map[string]string{
				"status": "deleted",
				"id":     id,
			})
			return
		}
	}

	http.Error(w, "playbook not found", http.StatusNotFound)
}