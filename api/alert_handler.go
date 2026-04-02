package api

import (
	"context"
	"encoding/json"
	"net/http"

	"sentinelx/models"
	"sentinelx/service"
)

type AlertHandler struct {
	Service *service.AlertService
}

func NewAlertHandler(s *service.AlertService) *AlertHandler {
	return &AlertHandler{Service: s}
}

func (h *AlertHandler) CreateAlert(w http.ResponseWriter, r *http.Request) {
	var alert models.Alert

	err := json.NewDecoder(r.Body).Decode(&alert)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// 🔥 call service
	h.Service.ProcessAlert(context.Background(), alert)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "alert processed",
	})
}