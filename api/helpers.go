package api

import (
	"encoding/json"
	"net/http"
)

// Standard error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// Write JSON response
func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// Write error response
func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, ErrorResponse{
		Error: message,
	})
}