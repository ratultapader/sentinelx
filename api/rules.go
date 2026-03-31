package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"strconv"
)

var rules = []map[string]interface{}{
	{"id": "1", "name": "SQL Detect", "condition": "req > 50", "action": "alert", "enabled": true},
	{"id": "2", "name": "Brute Force", "condition": "fail > 5", "action": "block", "enabled": false},
}

// ===============================
// GET RULES
// ===============================
func GetRules(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rules)
}

// ===============================
// CREATE RULE
// ===============================
func CreateRule(w http.ResponseWriter, r *http.Request) {
	var rule map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&rule)
	if err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	rule["id"] = strconv.Itoa(len(rules) + 1)
	rule["enabled"] = true

	println("DEBUG: creating rule", rule["id"])

	rules = append(rules, rule)

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(rule)
	if err != nil {
		http.Error(w, "failed to encode response", 500)
		return
	}
}

// ===============================
// TOGGLE RULE
// ===============================
func ToggleRule(w http.ResponseWriter, r *http.Request) {

	// Extract ID from URL
	path := strings.TrimPrefix(r.URL.Path, "/api/rules/")
	id := strings.TrimSuffix(path, "/toggle")

	for i, rule := range rules {
		if rule["id"] == id {
			rules[i]["enabled"] = !rule["enabled"].(bool)
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "toggled",
	})
}
