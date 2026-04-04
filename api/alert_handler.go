package api

import (
	// "context"
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

	ctx := r.Context() // ✅ use request context

	// ===============================
	// 🔹 Decode request
	// ===============================
	err := json.NewDecoder(r.Body).Decode(&alert)
	if err != nil {
		configs.ErrorCount.Inc()
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
	// 🔥 BLOCK CHECK (SAFE)
	// ===============================
	if ip != "" && storage.RDB != nil {
		val, err := storage.RDB.Get(ctx, "blocked:"+ip).Result()
		if err == nil && val == "1" {
			log.Println("🚫 REJECTED BLOCKED IP:", ip)
			configs.ErrorCount.Inc()
			http.Error(w, "IP is blocked", http.StatusForbidden)
			return
		}
	}

	// ===============================
	// 🔥 RATE LIMIT (SAFE)
	// ===============================
	if ip != "" && storage.RDB != nil {
		key := "rate_limit:" + ip

		count, err := storage.RDB.Incr(ctx, key).Result()
		if err == nil {
			storage.RDB.Expire(ctx, key, 60*time.Second)

			if count > 10 {
				log.Println("🚫 BLOCKED IP:", ip)

				storage.RDB.Set(ctx, "blocked:"+ip, "1", 5*time.Minute)

				configs.ErrorCount.Inc()
				http.Error(w, "IP blocked due to abuse", http.StatusTooManyRequests)
				return
			}
		}
	}

	// ===============================
	// 🔥 DEBUG
	// ===============================
	log.Println("🔥 ALERT RECEIVED:", alert)

	// ===============================
	// 🔹 Process alert
	// ===============================
	err = h.Service.ProcessAlert(ctx, alert)
	if err != nil {
		configs.ErrorCount.Inc()
		log.Println("❌ PROCESS ERROR:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("✅ ALERT PROCESSED")

	// ===============================
	// 🔥 CACHE INVALIDATION (SAFE)
	// ===============================
	if storage.RDB != nil {
		cacheKey := "incidents:" + tenantID
		err := storage.RDB.Del(ctx, cacheKey).Err()
		if err != nil {
			log.Println("⚠️ Failed to clear cache:", err)
		} else {
			log.Println("🧹 Cache cleared:", cacheKey)
		}
	}

	// ===============================
	// 🔥 REDIS PUB/SUB (SAFE)
	// ===============================
	if storage.RDB != nil {
		eventData, err := json.Marshal(alert)
		if err == nil {
			err = storage.RDB.Publish(ctx, "alerts_channel", eventData).Err()
			if err != nil {
				log.Println("⚠️ Failed to publish alert:", err)
			} else {
				log.Println("📡 Alert published to Redis")
			}
		}
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