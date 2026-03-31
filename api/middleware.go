package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

//////////////////////////////////////////////////////
// 🔥 RECOVERY MIDDLEWARE
//////////////////////////////////////////////////////

func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("PANIC:", err)
				debug.PrintStack()
				WriteError(w, http.StatusInternalServerError, "internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

//////////////////////////////////////////////////////
// 🔥 LOGGING MIDDLEWARE
//////////////////////////////////////////////////////

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("REQUEST:", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

//////////////////////////////////////////////////////
// 🔥 CORS MIDDLEWARE (FIXED)
//////////////////////////////////////////////////////

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// ✅ allow frontend origin (use "*" for dev)
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// ✅ IMPORTANT: include DELETE
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// ✅ allow headers
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Tenant-ID")

		// 🔥 handle preflight request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

//////////////////////////////////////////////////////
// 🔥 TENANT MIDDLEWARE (FIXED)
//////////////////////////////////////////////////////

func TenantMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// ✅ allow preflight without tenant
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		tenantID := r.Header.Get("X-Tenant-ID")
		fmt.Println("DEBUG MIDDLEWARE >>> HEADER TENANT:", tenantID)

		if tenantID == "" {
			tenantID = r.URL.Query().Get("tenant_id")
		}

		fmt.Println("DEBUG MIDDLEWARE >>> FINAL TENANT:", tenantID)

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