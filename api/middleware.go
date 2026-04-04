package api

import (
	"context"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/google/uuid"
	"sentinelx/configs"
)

//////////////////////////////////////////////////////
// 🔥 RECOVERY MIDDLEWARE
//////////////////////////////////////////////////////

func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {

				configs.Log("ERROR", "panic occurred", map[string]interface{}{
					"error": string(debug.Stack()),
				})

				WriteError(w, http.StatusInternalServerError, "internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

//////////////////////////////////////////////////////
// 🔥 REQUEST ID MIDDLEWARE (NEW)
//////////////////////////////////////////////////////

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		reqID := uuid.New().String()

		w.Header().Set("X-Request-ID", reqID)

		ctx := context.WithValue(r.Context(), "request_id", reqID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

//////////////////////////////////////////////////////
// 🔥 LOGGING MIDDLEWARE (UPGRADED)
//////////////////////////////////////////////////////

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start).Seconds()

		configs.RequestCount.WithLabelValues(r.Method, r.URL.Path).Inc()
configs.RequestDuration.Observe(duration)

		reqID, _ := r.Context().Value("request_id").(string)

		configs.Log("INFO", "request completed", map[string]interface{}{
			"method":     r.Method,
			"path":       r.URL.Path,
			"duration": duration,
			"request_id": reqID,
		})
	})
}

//////////////////////////////////////////////////////
// 🔥 CORS MIDDLEWARE (UNCHANGED)
//////////////////////////////////////////////////////

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Tenant-ID")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

//////////////////////////////////////////////////////
// 🔥 TENANT MIDDLEWARE (KEEP)
//////////////////////////////////////////////////////

func TenantMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// ✅ BYPASS FOR METRICS + HEALTH
if r.URL.Path == "/metrics" || r.URL.Path == "/health" {
	next.ServeHTTP(w, r)
	return
}

		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		tenantID := r.Header.Get("X-Tenant-ID")

		if tenantID == "" {
			tenantID = r.URL.Query().Get("tenant_id")
		}

		if tenantID == "" {
			http.Error(w, "missing tenant id", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), "tenant_id", tenantID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

//////////////////////////////////////////////////////
// 🔥 CHAIN BUILDER
//////////////////////////////////////////////////////

func Chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}