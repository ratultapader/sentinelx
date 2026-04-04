package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"sentinelx/configs"
	"sentinelx/models"
	"sentinelx/service"
	"sentinelx/storage"
)

type AlertHandler struct {
	Service *service.AlertService
}

func NewAlertHandler(s *service.AlertService) *AlertHandler {
	return &AlertHandler{Service: s}
}

func (h *AlertHandler) CreateAlert(w http.ResponseWriter, r *http.Request) {
	var alert models.Alert

	// ===============================
	// 🔹 Decode request
	// ===============================
	err := json.NewDecoder(r.Body).Decode(&alert)
	if err != nil {
		configs.ErrorCount.Inc() // ✅ metrics
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// ===============================
	// 🔹 Extract tenant
	// ===============================
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		tenantID = "default"
	}
	alert.TenantID = tenantID

	ip := alert.SourceIP

	// ===============================
	// 🔥 BLOCK CHECK
	// ===============================
	if ip != "" {
		val, _ := storage.RDB.Get(context.Background(), "blocked:"+ip).Result()
		if val == "1" {
			log.Println("🚫 REJECTED BLOCKED IP:", ip)
			configs.ErrorCount.Inc() // ✅ metrics
			http.Error(w, "IP is blocked", http.StatusForbidden)
			return
		}
	}

	// ===============================
	// 🔥 RATE LIMIT
	// ===============================
	if ip != "" {
		key := "rate_limit:" + ip

		count, _ := storage.RDB.Incr(context.Background(), key).Result()

		// expire in 60 seconds
		storage.RDB.Expire(context.Background(), key, 60*time.Second)

		if count > 10 {
			log.Println("🚫 BLOCKED IP:", ip)

			// 🔥 STORE BLOCKED IP
			storage.RDB.Set(context.Background(), "blocked:"+ip, "1", 5*time.Minute)

			configs.ErrorCount.Inc() // ✅ metrics
			http.Error(w, "IP blocked due to abuse", http.StatusTooManyRequests)
			return
		}
	}

	// ===============================
	// 🔥 DEBUG
	// ===============================
	log.Println("🔥 ALERT RECEIVED:", alert)

	// ===============================
	// 🔹 Process alert
	// ===============================
	err = h.Service.ProcessAlert(context.Background(), alert)
	if err != nil {
		configs.ErrorCount.Inc() // ✅ metrics
		log.Println("❌ PROCESS ERROR:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("✅ ALERT PROCESSED")

	// ===============================
	// 🔥 CACHE INVALIDATION
	// ===============================
	cacheKey := "incidents:" + tenantID
	err = storage.RDB.Del(context.Background(), cacheKey).Err()
	if err != nil {
		log.Println("⚠️ Failed to clear cache:", err)
	} else {
		log.Println("🧹 Cache cleared:", cacheKey)
	}

	// ===============================
	// 🔥 REDIS PUB/SUB
	// ===============================
	eventData, _ := json.Marshal(alert)

	err = storage.RDB.Publish(context.Background(), "alerts_channel", eventData).Err()
	if err != nil {
		log.Println("⚠️ Failed to publish alert:", err)
	} else {
		log.Println("📡 Alert published to Redis")
	}

	// ===============================
	// 🔹 Response
	// ===============================
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]string{
		"status": "alert processed",
	})
}